package app

import (
	"fmt"
	"net/http"
	_ "social-network-service/docs"
	"social-network-service/internal/api/account"
	"social-network-service/internal/api/dialog"
	"social-network-service/internal/api/post"
	"social-network-service/internal/api/user"
	"social-network-service/internal/config"
	"social-network-service/internal/database"
	"social-network-service/internal/middleware"
	"social-network-service/internal/repository"
	"social-network-service/internal/service"

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
	cfg := config.LoadConfig()

	cf := createConnectionFactory(cfg)

	tm := database.NewTransactionManager(cf)

	userRepository := repository.NewUserRepository(cf)
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
}

func createConnectionFactory(cfg *config.Config) *database.ConnectionFactory {
	dbCfg := cfg.Database

	cfCfg := database.ConnectionFactoryConfig{
		MasterConnectionString: fmt.Sprintf(
			"host=%v port=%v user=%v password=%v dbname=%v",
			dbCfg.Host,
			dbCfg.MasterPort,
			dbCfg.User,
			dbCfg.Password,
			dbCfg.DatabaseName),

		SyncConnectionString: fmt.Sprintf(
			"host=%v port=%v user=%v password=%v dbname=%v",
			dbCfg.Host,
			dbCfg.SyncPort,
			dbCfg.User,
			dbCfg.Password,
			dbCfg.DatabaseName),

		AsyncConnectionString: fmt.Sprintf(
			"host=%v port=%v user=%v password=%v dbname=%v",
			dbCfg.Host,
			dbCfg.AsyncPort,
			dbCfg.User,
			dbCfg.Password,
			dbCfg.DatabaseName),
	}

	return database.NewConnectionFactory(cfCfg)
}
