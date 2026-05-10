package helpers

import (
	"context"
	"database/sql"
)

type Transactor interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type sqlTransactor struct {
	db *sql.DB
}
type txKey struct{}

func NewTransactor(db *sql.DB) *sqlTransactor {
	return &sqlTransactor{db: db}
}

func (t *sqlTransactor) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	txCtx := context.WithValue(ctx, txKey{}, tx)
	if err := fn(txCtx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
func GetExecutor(ctx context.Context, db *sql.DB) DBExecutor {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return db
}
