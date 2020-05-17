package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/allisson/go-env"
	"github.com/allisson/go-pglock"
	"github.com/allisson/hammer"
	pb "github.com/allisson/hammer/api/v1"
	hammerGrpc "github.com/allisson/hammer/grpc"
	repository "github.com/allisson/hammer/repository/postgres"
	"github.com/allisson/hammer/service"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	logger       *zap.Logger
	sqlDB        *sqlx.DB
	grpcEndpoint string
	httpEndpoint string
)

type appContext struct {
	topicService           hammer.TopicService
	subscriptionService    hammer.SubscriptionService
	messageService         hammer.MessageService
	deliveryService        hammer.DeliveryService
	deliveryAttemptService hammer.DeliveryAttemptService
	migrationService       hammer.MigrationService
}

func newAppContext() appContext {
	// Create repositories
	topicRepo := repository.NewTopic(sqlDB)
	subscriptionRepo := repository.NewSubscription(sqlDB)
	messageRepo := repository.NewMessage(sqlDB)
	deliveryRepo := repository.NewDelivery(sqlDB)
	deliveryAttemptRepo := repository.NewDeliveryAttempt(sqlDB)
	txFactoryRepo := repository.NewTxFactory(sqlDB)
	migrationRepo := repository.NewMigration(sqlDB, env.GetString("HAMMER_DATABASE_MIGRATION_DIR", "file:///db/migrations"))

	// Create services
	topicService := service.NewTopic(&topicRepo, &txFactoryRepo)
	subscriptionService := service.NewSubscription(&topicRepo, &subscriptionRepo, &txFactoryRepo)
	messageService := service.NewMessage(&topicRepo, &messageRepo, &subscriptionRepo, &deliveryRepo, &txFactoryRepo)
	deliveryService := service.NewDelivery(&deliveryRepo, &deliveryAttemptRepo, &txFactoryRepo)
	deliveryAttemptService := service.NewDeliveryAttempt(&deliveryAttemptRepo)
	migrationService := service.NewMigration(&migrationRepo)

	return appContext{
		topicService:           &topicService,
		subscriptionService:    &subscriptionService,
		messageService:         &messageService,
		deliveryService:        &deliveryService,
		deliveryAttemptService: &deliveryAttemptService,
		migrationService:       &migrationService,
	}
}

func gatewayServer() {
	gatewayEnabled := env.GetBool("HAMMER_REST_API_ENABLED", true)
	if !gatewayEnabled {
		return
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterHammerHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		logger.Error("gateway-http-server", zap.Error(err))
		return
	}

	if err := http.ListenAndServe(httpEndpoint, mux); err != nil {
		logger.Error("gateway-http-server", zap.Error(err))
	}
}

func metricsServer() {
	metricsEnabled := env.GetBool("HAMMER_METRICS_ENABLED", true)
	if !metricsEnabled {
		return
	}
	port := env.GetInt("HAMMER_METRICS_PORT", 4001)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		logger.Error("metrics-server-failed-to-start", zap.Error(err))
	}
}

func init() {
	// Set logger
	logger, _ = zap.NewProduction()

	// Set database connection
	db, err := sqlx.Open("postgres", env.GetString("HAMMER_DATABASE_URL", ""))
	if err != nil {
		logger.Fatal("failed-to-start-database-client", zap.Error(err))
	}
	err = db.Ping()
	if err != nil {
		logger.Fatal("failed-to-ping-database", zap.Error(err))
	}
	sqlDB = db

	// Set grpc endpoint
	grpcEndpoint = fmt.Sprintf(":%d", env.GetInt("HAMMER_GRPC_PORT", 50051))

	// Set http endpoint
	httpEndpoint = fmt.Sprintf(":%d", env.GetInt("HAMMER_HTTP_PORT", 8000))
}

func main() {
	ac := newAppContext()
	app := cli.NewApp()
	app.Name = "Hammer"
	app.Usage = "CLI"
	app.Authors = []*cli.Author{
		{
			Name:  "Allisson Azevedo",
			Email: "allisson@gmail.com",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "Starts the server",
			Action: func(c *cli.Context) error {
				// Create grpc handlers
				topicHandler := hammerGrpc.NewTopicHandler(ac.topicService)
				subscriptionHandler := hammerGrpc.NewSubscriptionHandler(ac.subscriptionService)
				messageHandler := hammerGrpc.NewMessageHandler(ac.messageService)
				deliveryHandler := hammerGrpc.NewDeliveryHandler(ac.deliveryService)
				deliveryAttemptHandler := hammerGrpc.NewDeliveryAttemptHandler(ac.deliveryAttemptService)

				// Start tcp server
				listener, err := net.Listen("tcp", grpcEndpoint)
				if err != nil {
					logger.Fatal("failed-to-listen", zap.Error(err))
				}

				// Start http gateway
				go gatewayServer()

				// Create grpc server
				grpcServer := grpc.NewServer(
					grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
						grpc_ctxtags.StreamServerInterceptor(),
						grpc_prometheus.StreamServerInterceptor,
						grpc_zap.StreamServerInterceptor(logger),
						grpc_recovery.StreamServerInterceptor(),
					)),
					grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
						grpc_ctxtags.UnaryServerInterceptor(),
						grpc_prometheus.UnaryServerInterceptor,
						grpc_zap.UnaryServerInterceptor(logger),
						grpc_recovery.UnaryServerInterceptor(),
					)),
				)
				server := hammerGrpc.NewServer(topicHandler, subscriptionHandler, messageHandler, deliveryHandler, deliveryAttemptHandler)

				// Register grpc services
				pb.RegisterHammerServer(grpcServer, &server)

				// Enable grpc_prometheus histograms
				grpc_prometheus.EnableHandlingTimeHistogram()

				// Register grpc_prometheus default metrics
				grpc_prometheus.Register(grpcServer)

				// Load metrics server
				go metricsServer()

				// Start grpc server and make graceful shutdown
				idleConnsClosed := make(chan struct{})
				go func() {
					sigint := make(chan os.Signal, 1)

					// interrupt signal sent from terminal
					signal.Notify(sigint, os.Interrupt)
					// sigterm signal sent from kubernetes
					signal.Notify(sigint, syscall.SIGTERM)

					<-sigint

					// We received an interrupt signal, shut down.
					logger.Info("grpc-server-shutdown-started")
					grpcServer.GracefulStop()
					close(idleConnsClosed)
					logger.Info("grpc-server-shutdown-finished")
				}()

				logger.Info("grpc-server-started")
				if err := grpcServer.Serve(listener); err != nil {
					logger.Error("grpc-server-serve", zap.Error(err))
				}

				<-idleConnsClosed

				return nil
			},
		},
		{
			Name:    "migrate",
			Aliases: []string{"m"},
			Usage:   "Run dabatase migrate",
			Action: func(c *cli.Context) error {
				err := ac.migrationService.Run()
				if err != nil {
					return err
				}
				logger.Info("database-migrations-completed")
				return nil
			},
		},
		{
			Name:    "worker",
			Aliases: []string{"w"},
			Usage:   "Starts the worker",
			Action: func(c *cli.Context) error {
				// Create lock
				sqlConn, err := sqlDB.DB.Conn(context.Background())
				if err != nil {
					return err
				}
				lock := pglock.NewLock(sqlConn)

				// Create worker service
				workerService := service.NewWorker(&lock, ac.deliveryService)

				idleConnsClosed := make(chan struct{})
				go func() {
					sigint := make(chan os.Signal, 1)

					// interrupt signal sent from terminal
					signal.Notify(sigint, os.Interrupt)
					// sigterm signal sent from kubernetes
					signal.Notify(sigint, syscall.SIGTERM)

					<-sigint

					// We received an interrupt signal, shut down.
					logger.Info("worker-shutdown-started")
					if err := workerService.Stop(); err != nil {
						logger.Error("worker-service-stop", zap.Error(err))
					}
					close(idleConnsClosed)
					logger.Info("worker-shutdown-finished")
				}()

				logger.Info("worker-started")
				if err := workerService.Run(); err != nil {
					logger.Error("worker-service-run", zap.Error(err))
				}

				<-idleConnsClosed

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal("app", zap.Error(err))
	}
}
