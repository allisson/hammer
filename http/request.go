package http

import (
	"net/http"
	"strconv"

	"github.com/allisson/hammer"
)

func getLimitOffset(r *http.Request) (int, int) {
	limit := hammer.DefaultPaginationLimit
	offset := 0

	queryLimit := r.URL.Query().Get("limit")
	queryOffset := r.URL.Query().Get("offset")
	if queryLimit != "" {
		l, err := strconv.Atoi(queryLimit)
		if err == nil && l > 0 {
			limit = l
		}
	}
	if queryOffset != "" {
		o, err := strconv.Atoi(queryOffset)
		if err == nil && o >= 0 {
			offset = o
		}
	}

	return limit, offset
}
