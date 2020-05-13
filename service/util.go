package service

import (
	"math/rand"
	mathrand "math/rand"
	"time"

	"github.com/allisson/hammer"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
)

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func generateULID() (string, error) {
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
