package http

import (
	"net/http"
)

// PingHandler implements Get method for connectivity test
type PingHandler struct{}

// Get returns status code 204
func (p *PingHandler) Get(w http.ResponseWriter, r *http.Request) {
	contentType := "application/json"
	makeResponse(w, []byte{}, http.StatusNoContent, contentType)

}

// NewPingHandler returns a new PingHandler
func NewPingHandler() PingHandler {
	return PingHandler{}
}
