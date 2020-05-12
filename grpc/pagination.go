package grpc

import "github.com/allisson/hammer"

func parsePagination(limit, offset uint32) (uint, uint) {
	if limit == 0 {
		limit = uint32(hammer.DefaultPaginationLimit)
	} else if limit > uint32(hammer.MaxPaginationLimit) {
		limit = uint32(hammer.MaxPaginationLimit)
	}
	return uint(limit), uint(offset)
}
