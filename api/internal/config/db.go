package config

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	p, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	// test connection
	if err := p.Ping(ctx); err != nil {
		return nil, err
	}

	return p, nil
}
