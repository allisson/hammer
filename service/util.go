package service

import (
	mathrand "math/rand"
	"time"

	"github.com/allisson/hammer"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
)

func generateID() (string, error) {
	seed := time.Now().UnixNano()
	source := mathrand.NewSource(seed)
	entropy := mathrand.New(source)
	id, err := ulid.New(ulid.Timestamp(time.Now()), entropy)
	return id.String(), err
}

func rollback(tx hammer.TxRepository, message string) {
	err := tx.Rollback()
	if err != nil {
		logger.Error(message, zap.Error(err))
	}
}
