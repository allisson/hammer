package service

import (
	"testing"

	"github.com/allisson/hammer/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMigration(t *testing.T) {
	migrationRepo := &mocks.MigrationRepository{}
	migrationService := NewMigration(migrationRepo)
	migrationRepo.On("Run").Return(nil)

	err := migrationService.Run()
	assert.Nil(t, err)
}
