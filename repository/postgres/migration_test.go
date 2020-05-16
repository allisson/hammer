package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigration(t *testing.T) {
	th := newTxnTestHelper()
	defer th.db.Close()
	migrationDir := "file://../../db/migrations"
	migration := NewMigration(th.db, migrationDir)
	err := migration.Run()
	assert.Nil(t, err)
}
