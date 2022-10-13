package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/allisson/go-env"
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
	"google.golang.org/grpc/credentials/insecure"

	"github.com/allisson/hammer"
	pb "github.com/allisson/hammer/api/v1"
	hammerGrpc "github.com/allisson/hammer/grpc"
	repository "github.com/allisson/hammer/repository/postgres"
	"github.com/allisson/hammer/service"
)

var (
	sqlDB             *sqlx.DB
	sqlConn           *sql.Conn
	grpcEndpoint      string
	httpEndpoint      string
	readHeaderTimeout time.Duration
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
	migrationRepo := repository.NewMigration(sqlDB, env.GetString("HAMMER_DATABASE_MIGRATION_DIR", "file:///db/migrations"))

	// Create services
	topicService := service.NewTopic(topicRepo)
	subscriptionService := service.NewSubscription(topicRepo, subscriptionRepo)
	messageService := service.NewMessage(topicRepo, messageRepo, subscriptionRepo, deliveryRepo)
	deliveryService := service.NewDelivery(deliveryRepo, deliveryAttemptRepo)
	deliveryAttemptService := service.NewDeliveryAttempt(deliveryAttemptRepo)
	migrationService := service.NewMigration(migrationRepo)

	return appContext{
		topicService:           topicService,
		subscriptionService:    subscriptionService,
		messageService:         messageService,
		deliveryService:        deliveryService,
		deliveryAttemptService: deliveryAttemptService,
		migrationService:       migrationService,
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
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterHammerHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts)
	if err != nil {
		zap.L().Error("gateway-http-server", zap.Error(err))
		return
	}
	server := &http.Server{
		Addr:              httpEndpoint,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}
	if err := server.ListenAndServe(); err != nil {
		zap.L().Error("gateway-http-server", zap.Error(err))
	}
}

func metricsServer() {
	metricsEnabled := env.GetBool("HAMMER_METRICS_ENABLED", true)
	if !metricsEnabled {
		return
	}
	port := env.GetInt("HAMMER_METRICS_PORT", 4001)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}
	err := server.ListenAndServe()
	if err != nil {
		zap.L().Error("metrics-server-failed-to-start", zap.Error(err))
	}
}

func healthCheckServer() {
	healthCheckEnabled := env.GetBool("HAMMER_HEALTH_CHECK_ENABLED", true)
	if !healthCheckEnabled {
		return
	}
	port := env.GetInt("HAMMER_HEALTH_CHECK_PORT", 9000)
	mux := http.NewServeMux()
	handler := func(w http.ResponseWriter, r *http.Request) {
		if err := sqlDB.Ping(); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		if sqlConn != nil {
			if err := sqlConn.PingContext(context.Background()); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusNoContent)
	}
	mux.HandleFunc("/liveness", handler)
	mux.HandleFunc("/readiness", handler)
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}
	err := server.ListenAndServe()
	if err != nil {
		zap.L().Error("health-check-server-failed-to-start", zap.Error(err))
	}
}

func init() {
	// Set logger
	logger, _ := zap.NewProduction()
	_ = zap.ReplaceGlobals(logger)

	// Set database connection
	db, err := sqlx.Open("postgres", env.GetString("HAMMER_DATABASE_URL", ""))
	if err != nil {
		zap.L().Fatal("failed-to-start-database-client", zap.Error(err))
	}
	err = db.Ping()
	if err != nil {
		zap.L().Fatal("failed-to-ping-database", zap.Error(err))
	}
	db.SetMaxOpenConns(env.GetInt("HAMMER_DATABASE_MAX_OPEN_CONNS", 3))
	sqlDB = db

	// Set grpc endpoint
	grpcEndpoint = fmt.Sprintf(":%d", env.GetInt("HAMMER_GRPC_PORT", 50051))

	// Set http endpoint
	httpEndpoint = fmt.Sprintf(":%d", env.GetInt("HAMMER_HTTP_PORT", 8000))

	// Set ReadHeaderTimeout
	readHeaderTimeout = 60 * time.Second
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
					zap.L().Fatal("failed-to-listen", zap.Error(err))
				}

				// Start health check
				go healthCheckServer()

				// Start http gateway
				go gatewayServer()

				// Create grpc server
				grpcServer := grpc.NewServer(
					grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
						grpc_ctxtags.StreamServerInterceptor(),
						grpc_prometheus.StreamServerInterceptor,
						grpc_zap.StreamServerInterceptor(zap.L()),
						grpc_recovery.StreamServerInterceptor(),
					)),
					grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
						grpc_ctxtags.UnaryServerInterceptor(),
						grpc_prometheus.UnaryServerInterceptor,
						grpc_zap.UnaryServerInterceptor(zap.L()),
						grpc_recovery.UnaryServerInterceptor(),
					)),
				)
				server := hammerGrpc.NewServer(topicHandler, subscriptionHandler, messageHandler, deliveryHandler, deliveryAttemptHandler)

				// Register grpc services
				pb.RegisterHammerServer(grpcServer, &server)

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
					zap.L().Info("grpc-server-shutdown-started")
					grpcServer.GracefulStop()
					close(idleConnsClosed)
					zap.L().Info("grpc-server-shutdown-finished")
				}()

				zap.L().Info("grpc-server-started")
				if err := grpcServer.Serve(listener); err != nil {
					zap.L().Error("grpc-server-serve", zap.Error(err))
				}

				<-idleConnsClosed

				return nil
			},
		},
		{
			Name:    "migrate",
			Aliases: []string{"m"},
			Usage:   "Run database migrate",
			Action: func(c *cli.Context) error {
				err := ac.migrationService.Run(c.Context)
				if err != nil {
					return err
				}
				zap.L().Info("database-migrations-completed")
				return nil
			},
		},
		{
			Name:    "worker",
			Aliases: []string{"w"},
			Usage:   "Starts the worker",
			Action: func(c *cli.Context) error {
				// Create repositories
				deliveryRepo := repository.NewDelivery(sqlDB)

				// Create worker service
				pollingInterval := time.Duration(env.GetInt("HAMMER_WORKER_DATABASE_DELAY", 1)) * time.Second
				workerService := service.NewWorker(deliveryRepo, pollingInterval)

				// Start health check
				go healthCheckServer()

				// Start worker
				workerService.Run(c.Context)

				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		zap.L().Error("app", zap.Error(err))
	}
}
