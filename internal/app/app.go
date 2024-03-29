package app

import (
	"context"
	"database/sql"
	_ "social-network-service/docs"
	"social-network-service/internal/api/account"
	"social-network-service/internal/api/dialog"
	"social-network-service/internal/api/post"
	"social-network-service/internal/api/user"
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
	ctx := context.Background()

	connStr := "host=localhost port=5432 user=postgres password=123 dbname=social_network_db"

	db, err := sql.Open("pgx", connStr)

	if err != nil {
		panic(err)
	}

	err = db.PingContext(ctx)

	if err != nil {
		panic(err)
	}

	tm := NewTransactionManager(db)

	userRepository := repository.NewUserRepository(db)
	userAccountRepository := repository.NewUserAccountRepository(db)
	dialogRepository := repository.NewDialogRepository(db)
	postRepository := repository.NewPostRepository(db)
	userFriendRepository := repository.NewUserFriendRepository(db)

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

	// TODO: Improve working with endpoints.
	_ = account.RegisterAccountEndpoints(appService, engine)

	_ = user.RegisterUserOpenEndpoints(appService, engine)

	userRouter := user.RegisterUserClosedEndpoints(appService, jwtService, engine)
	userRouter.Use(authMiddleware)

	postRouter := post.RegisterPostEndpoints(appService, jwtService, engine)
	postRouter.Use(authMiddleware)

	dialogRouter := dialog.RegisterDialogEndpoints(appService, jwtService, engine)
	dialogRouter.Use(authMiddleware)

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	engine.Run(":8080")
}
