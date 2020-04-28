package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/allisson/hammer"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func init() {
	logger, _ = zap.NewProduction()
}

func makeResponse(w http.ResponseWriter, body []byte, statusCode int, contentType string) {
	w.Header().Set("Content-Type", fmt.Sprintf("%s; charset=utf-8", contentType))
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		logger.Error("makeResponse-failed-to-write-response-body", zap.Error(err))
	}
}

func errorResponse(w http.ResponseWriter, err hammer.Error, contentType string) {
	responseBody, _ := json.Marshal(&err)
	makeResponse(w, responseBody, err.Code, contentType)
}
