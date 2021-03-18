package service

import (
	"context"

	"github.com/allisson/hammer"
)

// Migration is a implementation of hammer.MigrationService
type Migration struct {
	migrationService hammer.MigrationService
}

// Run migrations
func (m Migration) Run(ctx context.Context) error {
	return m.migrationService.Run(ctx)
}

// NewMigration will create a implementation of hammer.MigrationService
func NewMigration(migrationService hammer.MigrationService) Migration {
	return Migration{migrationService: migrationService}
}
