package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	_ "social-network-service/docs"
	"social-network-service/internal/api/account"
	"social-network-service/internal/api/dialog"
	"social-network-service/internal/api/post"
	"social-network-service/internal/api/user"
	"social-network-service/internal/cache"
	"social-network-service/internal/config"
	"social-network-service/internal/database"
	"social-network-service/internal/kafka/consumer"
	"social-network-service/internal/kafka/producer"
	"social-network-service/internal/middleware"
	"social-network-service/internal/repository"
	"social-network-service/internal/service"
	"social-network-service/internal/ws"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securitydefinitions.apikey bearer
// @in header
// @name Authorization

func Run() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	cfg := config.LoadConfig()

	cf := createConnectionFactory(cfg)

	tm := database.NewTransactionManager(cf)

	userRepository := repository.NewUserRepository(cf)
	userAccountRepository := repository.NewUserAccountRepository(cf)
	userFriendRepository := repository.NewUserFriendRepository(cf)
	dialogRepository := repository.NewDialogRepository(cf)
	postRepository := repository.NewPostRepository(cf)

	jwtService := service.NewJwtService()

	feedCacheNotifier, err := producer.NewFeedCommandProducer(cfg.KafkaBrokers, cfg.Producers.Feed.Topic)

	if err != nil {
		panic(err)
	}

	postEventNotifier, err := producer.NewPostEventProducer(cfg.KafkaBrokers, cfg.Producers.Posts.Topic)

	if err != nil {
		panic(err)
	}

	redisOptions := &redis.Options{
		Addr: cfg.Redis.ConnectionString,
	}

	redisClient := redis.NewClient(redisOptions)

	feedCache := cache.NewFeedCache(redisClient)

	wsHub := ws.NewHub()
	userNotifier := ws.NewUserNotifier(redisClient, wsHub)

	appServiceConfig := &service.AppServiceConfiguration{
		TokenGenerator:        jwtService,
		UserRepository:        userRepository,
		UserAccountRepository: userAccountRepository,
		UserFriendRepository:  userFriendRepository,
		DialogRepository:      dialogRepository,
		PostRepository:        postRepository,
		FeedCache:             feedCache,
		FeedCacheNotifier:     feedCacheNotifier,
		PostEventNotifier:     postEventNotifier,
		UserNotifier:          userNotifier,
		TransactionManager:    tm,
	}

	appService := service.NewAppService(appServiceConfig)

	go wsHub.Run(ctx)
	go userNotifier.Subscribe(ctx)

	engine := gin.Default()

	errorHandlingMiddleware := middleware.NewErrorHandlingMiddleware()
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	engine.Use(errorHandlingMiddleware)

	account.RegisterAccountEndpoints(appService, engine)
	user.RegisterUserClosedEndpoints(appService, jwtService, engine, authMiddleware)
	post.RegisterPostEndpoints(appService, jwtService, engine, authMiddleware)
	dialog.RegisterDialogEndpoints(appService, jwtService, engine, authMiddleware)

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	wg := new(sync.WaitGroup)

	feecCommandConsumer := consumer.NewFeedCommandConsumer(appService)
	postEventConsumer := consumer.NewPostEventConsumer(appService)

	// TODO: Consider using only one consumer to consume messages from both topics.
	wg.Add(2)
	consumer.UseFeedCommandConsumer(ctx, cfg.KafkaBrokers, cfg.Consumers.Feed.Topic, feecCommandConsumer, wg)
	consumer.UsePostEventConsumer(ctx, cfg.KafkaBrokers, cfg.Consumers.Posts.Topic, postEventConsumer, wg)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	go func() {
		engine.Run(":8080")
	}()

	go func() {
		hubMux := http.NewServeMux()
		hubMux.Handle("/hub", ws.HandleHub(wsHub, jwtService))
		http.ListenAndServe(":8081", hubMux)
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)

	select {
	case <-sigint:
		log.Println("SIGTERM received, graceful shutdown started")
		cancel()
	case <-sigterm:
		log.Println("SIGTERM received, graceful shutdown started")
		cancel()
	}

	wg.Wait()

	log.Println("graceful shutdown finished, have a nice day")
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
