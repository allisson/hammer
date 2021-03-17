package repository

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// For migrate with files
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

// Migration is a implementation of hammer.MigrationRepository
type Migration struct {
	db           *sqlx.DB
	migrationDir string
}

// Run migrations
func (m *Migration) Run(ctx context.Context) error {
	driver, err := postgres.WithInstance(m.db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	mg, err := migrate.NewWithDatabaseInstance(m.migrationDir, "postgres", driver)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = mg.Up()
	if err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
	}

	return nil
}

// NewMigration will create a implementation of hammer.MigrationRepository
func NewMigration(db *sqlx.DB, migrationDir string) *Migration {
	return &Migration{
		db:           db,
		migrationDir: migrationDir,
	}
}
