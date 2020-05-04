package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/allisson/go-env"
	h "github.com/allisson/hammer/http"
	repository "github.com/allisson/hammer/repository/postgres"
	"github.com/allisson/hammer/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	sqlDB  *sqlx.DB
)

func init() {
	// Set logger
	logger, _ = zap.NewProduction()

	// Set database connection
	db, err := sqlx.Open("postgres", env.GetString("DATABASE_URL", ""))
	if err != nil {
		logger.Fatal("failed-to-start-database-client", zap.Error(err))
	}
	err = db.Ping()
	if err != nil {
		logger.Fatal("failed-to-ping-database", zap.Error(err))
	}
	sqlDB = db
}

func main() {
	// Start chi router
	r := chi.NewRouter()
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Create repositories
	topicRepo := repository.NewTopic(sqlDB)
	subscriptionRepo := repository.NewSubscription(sqlDB)
	messageRepo := repository.NewMessage(sqlDB)
	deliveryRepo := repository.NewDelivery(sqlDB)
	txFactoryRepo := repository.NewTxFactory(sqlDB)

	// Create services
	topicService := service.NewTopic(&topicRepo, &txFactoryRepo)
	subscriptionService := service.NewSubscription(&topicRepo, &subscriptionRepo, &txFactoryRepo)
	messageService := service.NewMessage(&topicRepo, &messageRepo, &subscriptionRepo, &deliveryRepo, &txFactoryRepo)

	// Create handlers
	pingHandler := h.NewPingHandler()
	topicHandler := h.NewTopicHandler(&topicService)
	subscriptionHandler := h.NewSubscriptionHandler(&subscriptionService)
	messageHandler := h.NewMessageHandler(&messageService)

	// Create routes
	r.Get("/ping", pingHandler.Get)
	r.Post("/topics", topicHandler.Create)
	r.Get("/topics", topicHandler.List)
	r.Post("/subscriptions", subscriptionHandler.Create)
	r.Get("/subscriptions", subscriptionHandler.List)
	r.Post("/messages", messageHandler.Create)
	r.Get("/messages", messageHandler.List)

	// Start server and make graceful shutdown
	logger.Info("start-http-server")
	port := fmt.Sprintf(":%s", env.GetString("PORT", "8000"))
	httpServer := &http.Server{Addr: port, Handler: r}
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)

		// interrupt signal sent from terminal
		signal.Notify(sigint, os.Interrupt)
		// sigterm signal sent from kubernetes
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		// We received an interrupt signal, shut down.
		logger.Info("http-server-shutdown-started")
		if err := httpServer.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			logger.Error("http-server-shutdown", zap.Error(err))
		}
		close(idleConnsClosed)
		logger.Info("http-server-shutdown-finished")
	}()

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		logger.Error("http-server-listen-and-server", zap.Error(err))
	}

	<-idleConnsClosed
}