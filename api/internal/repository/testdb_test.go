package repository

import (
	"cafenetchi-api/internal/config"
	db "cafenetchi-api/internal/db/generated"
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

type testRepo struct {
	User User
	tx   pgx.Tx
	ctx  context.Context
}

func setupUserRepo(t *testing.T) *testRepo {
	t.Helper()
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		t.Fatal("DB_URL is not set")
	}
	ctx := context.Background()
	pool, err := config.NewPool(ctx, dsn)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		pool.Close()
		t.Fatalf("begin transaction: %v", err)
	}
	t.Cleanup(func() {
		_ = tx.Rollback(ctx)
		pool.Close()
	})

	queries := db.New(tx)
	return &testRepo{
		User: NewUser(queries),
		tx:   tx,
		ctx:  ctx,
	}

}
