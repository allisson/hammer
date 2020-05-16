package service

import "github.com/allisson/hammer"

// Migration is a implementation of hammer.MigrationService
type Migration struct {
	migrationService hammer.MigrationService
}

// Run migrations
func (m *Migration) Run() error {
	return m.migrationService.Run()
}

// NewMigration will create a implementation of hammer.MigrationService
func NewMigration(migrationService hammer.MigrationService) Migration {
	return Migration{migrationService: migrationService}
}
