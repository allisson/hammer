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
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func init() {
	l, _ := zap.NewProduction()
	logger = l
}

func main() {
	// Start chi router
	r := chi.NewRouter()
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Add handlers
	pingHandler := h.NewPingHandler()
	r.Get("/ping", pingHandler.Get)

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
