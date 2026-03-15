package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Open(ctx context.Context, pgDatabaseUrl string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(pgDatabaseUrl)
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = 10
	cfg.MinIdleConns = 0
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.MaxConnLifetime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// check ping
	ctxPing, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := pool.Ping(ctxPing); err != nil {
		return nil, err
	}

	return pool, nil
}
