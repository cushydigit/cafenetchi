package config

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func InitPool(ctx context.Context, dsn string) error {
	p, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return err
	}

	if err := p.Ping(ctx); err != nil {
		log.Fatal("Failed to ping to database:", err)
		return err
	}

	log.Println("✅ Database connected successfully")

	pool = p
	return nil
}

func GetPool(ctx context.Context) *pgxpool.Pool {
	if pool == nil {
		log.Println("try to ref to nil object")
		return nil
	}
	return pool
}

func Close() error {
	if pool != nil {
		pool.Close()
	}
	return nil
}
