package http

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func init() {
	logger, _ = zap.NewProduction()
}

// MakeResponse write response on ResponseWriter
func MakeResponse(w http.ResponseWriter, body []byte, statusCode int, contentType string) {
	w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", contentType))
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		logger.Error("makeResponse-failed-to-write-response-body", zap.Error(err))
	}
}
