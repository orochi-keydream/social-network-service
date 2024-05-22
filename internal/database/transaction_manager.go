package database

import (
	"context"
	"database/sql"
)

type TransactionManager struct {
	cf *ConnectionFactory
}

func NewTransactionManager(cf *ConnectionFactory) *TransactionManager {
	return &TransactionManager{
		cf: cf,
	}
}

func (tm *TransactionManager) Begin(ctx context.Context) (*sql.Tx, error) {
	opts := sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	}

	tx, err := tm.cf.GetMaster().BeginTx(ctx, &opts)

	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (tm *TransactionManager) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (tm *TransactionManager) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}
