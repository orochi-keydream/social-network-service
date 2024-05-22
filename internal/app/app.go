package app

import (
	_ "social-network-service/docs"
	"social-network-service/internal/api/account"
	"social-network-service/internal/api/dialog"
	"social-network-service/internal/api/post"
	"social-network-service/internal/api/user"
	"social-network-service/internal/database"
	"social-network-service/internal/middleware"
	"social-network-service/internal/repository"
	"social-network-service/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securitydefinitions.apikey bearer
// @in header
// @name Authorization

func Run() {
	cfCfg := database.ConnectionFactoryConfig{
		MasterConnectionString: "host=haproxy port=15432 user=postgres password=123 dbname=social_network_db",
		SyncConnectionString:   "host=haproxy port=25432 user=postgres password=123 dbname=social_network_db",
		AsyncConnectionString:  "host=haproxy port=35432 user=postgres password=123 dbname=social_network_db",
	}

	cf := database.NewConnectionFactory(cfCfg)

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

	engine.Run(":8080")
}
