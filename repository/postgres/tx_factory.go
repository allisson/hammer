package repository

import (
	"github.com/allisson/hammer"
	"github.com/jmoiron/sqlx"
)

// TxFactory is a implementation of hammer.TxFactoryRepository
type TxFactory struct {
	db *sqlx.DB
}

// New returns a hammer.TxRepository
func (t *TxFactory) New() (hammer.TxRepository, error) {
	sqlTx, err := t.db.Beginx()
	if err != nil {
		return nil, err
	}
	tx := NewTx(sqlTx)
	return &tx, nil
}

// NewTxFactory returns a new TxFactory with db connection
func NewTxFactory(db *sqlx.DB) TxFactory {
	return TxFactory{db: db}
}
