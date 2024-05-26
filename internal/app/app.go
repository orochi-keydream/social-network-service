package app

import (
	"fmt"
	"net/http"
	"os"
	_ "social-network-service/docs"
	"social-network-service/internal/api/account"
	"social-network-service/internal/api/dialog"
	"social-network-service/internal/api/post"
	"social-network-service/internal/api/user"
	"social-network-service/internal/database"
	"social-network-service/internal/middleware"
	"social-network-service/internal/repository"
	"social-network-service/internal/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securitydefinitions.apikey bearer
// @in header
// @name Authorization

func Run() {
	dbHost := getDatabaseHost()

	cfCfg := database.ConnectionFactoryConfig{
		MasterConnectionString: fmt.Sprintf("host=%s port=15432 user=postgres password=123 dbname=social_network_db", dbHost),
		SyncConnectionString:   fmt.Sprintf("host=%s port=25432 user=postgres password=123 dbname=social_network_db", dbHost),
		AsyncConnectionString:  fmt.Sprintf("host=%s port=35432 user=postgres password=123 dbname=social_network_db", dbHost),
	}

	cf := database.NewConnectionFactory(cfCfg)

	tm := database.NewTransactionManager(cf)

	userRepositoryConfig := repository.UserRepositoryConfiguartion{
		UseAsyncReplicaForReadOperations: shouldUseAsyncReplica(),
	}

	userRepository := repository.NewUserRepository(userRepositoryConfig, cf)

	userAccountRepository := repository.NewUserAccountRepository(cf)
	dialogRepository := repository.NewDialogRepository(cf)
	postRepository := repository.NewPostRepository(cf)
	userFriendRepository := repository.NewUserFriendRepository(cf)

	jwtService := service.NewJwtService()

	appServiceConfig := &service.AppServiceConfiguration{
		TokenGenerator:        jwtService,
		UserRepository:        userRepository,
		UserAccountRepository: userAccountRepository,
		UserFriendRepository:  userFriendRepository,
		DialogRepository:      dialogRepository,
		PostRepository:        postRepository,
		TransactionManager:    tm,
	}

	appService := service.NewAppService(appServiceConfig)

	engine := gin.Default()

	errorHandlingMiddleware := middleware.NewErrorHandlingMiddleware()
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	engine.Use(errorHandlingMiddleware)

	account.RegisterAccountEndpoints(appService, engine)
	user.RegisterUserClosedEndpoints(appService, jwtService, engine, authMiddleware)
	post.RegisterPostEndpoints(appService, jwtService, engine, authMiddleware)
	dialog.RegisterDialogEndpoints(appService, jwtService, engine, authMiddleware)

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	engine.Run(":8080")
	time.Sleep(time.Second * 5)
}

func shouldUseAsyncReplica() bool {
	str := os.Getenv("USE_ASYNC_REPLICA")

	if str == "" {
		return false
	}

	val, err := strconv.ParseBool(str)

	if err != nil {
		panic(err)
	}

	return val
}

func getDatabaseHost() string {
	isRunningInContainerStr := os.Getenv("IS_RUNNING_IN_CONTAINER")

	var isRunningInContainer bool

	if isRunningInContainerStr != "" {
		var err error

		isRunningInContainer, err = strconv.ParseBool(isRunningInContainerStr)

		if err != nil {
			panic(err)
		}
	}

	if isRunningInContainer {
		return "haproxy"
	} else {
		return "localhost"
	}
}
