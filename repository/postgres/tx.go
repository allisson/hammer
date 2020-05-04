package repository

import "github.com/jmoiron/sqlx"

// Tx is a implementation of hammer.TxRepository
type Tx struct {
	tx *sqlx.Tx
}

// Exec executes a query that doesn't return rows
func (t *Tx) Exec(query string, arg interface{}) error {
	_, err := t.tx.NamedExec(query, arg)
	return err
}

// Commit commits the transaction
func (t *Tx) Commit() error {
	return t.tx.Commit()
}

// Rollback aborts the transaction
func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}

// NewTx returns a new Tx with *sql.Tx
func NewTx(tx *sqlx.Tx) Tx {
	return Tx{tx: tx}
}
