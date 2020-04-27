package service

import (
	mathrand "math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func generateID() (string, error) {
	seed := time.Now().UnixNano()
	source := mathrand.NewSource(seed)
	entropy := mathrand.New(source)
	id, err := ulid.New(ulid.Timestamp(time.Now()), entropy)
	return id.String(), err
}
