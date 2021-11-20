package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/allisson/hammer/mocks"
)

func TestMigration(t *testing.T) {
	ctx := context.Background()
	migrationRepo := &mocks.MigrationRepository{}
	migrationService := NewMigration(migrationRepo)
	migrationRepo.On("Run", mock.Anything).Return(nil)

	err := migrationService.Run(ctx)
	assert.Nil(t, err)
}
