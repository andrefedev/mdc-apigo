package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ctxKeyTx = struct{ name string }{"ctx-key-tx"}

func txFromContext(ctx context.Context) (*pgxpool.Tx, bool) {
	tx, ok := ctx.Value(ctxKeyTx).(*pgxpool.Tx)
	return tx, ok
}

type Pgdb struct {
	pool *pgxpool.Pool
}

func NewPgdb(pool *pgxpool.Pool) *Pgdb {
	return &Pgdb{pool: pool}
}

// Exec executes a query without returning any rows.
func (x Pgdb) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	if tx, ok := txFromContext(ctx); ok {
		return tx.Exec(ctx, query, args...)
	}

	return x.pool.Exec(ctx, query, args...)
}

// Query executes a query that returns rows, typically a SELECT.
func (x Pgdb) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	if tx, ok := txFromContext(ctx); ok {
		return tx.Query(ctx, query, args...)
	}
	return x.pool.Query(ctx, query, args...)
}

// QueryRow executes a query expected to return at most one row.
func (x Pgdb) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	if tx, ok := txFromContext(ctx); ok {
		return tx.QueryRow(ctx, query, args...)
	}
	return x.pool.QueryRow(ctx, query, args...)
}

// WithTx will start a new SQL transaction and hold a reference
// to the transaction inside the context.
// Next calls within the txFunc will use the new transaction from the context.
func (x Pgdb) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, ok := txFromContext(ctx); ok {
		return fn(ctx)
	}

	// Start a transaction.
	tx, err := x.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback(ctx)

	newContext := context.WithValue(ctx, ctxKeyTx, tx)
	if err := fn(newContext); err != nil {
		return fmt.Errorf("Pgdb.WithTx: fn(): [%w]", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("Pgdb.WithTx: tx.Commit(): [%w]", err)
	}

	return nil
}
