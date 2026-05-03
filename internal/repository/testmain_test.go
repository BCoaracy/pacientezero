//go:build integration

package repository_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=pacientezero password=pacientezero_dev_123 dbname=pacientezero_db sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		panic("falha ao conectar no banco de teste: " + err.Error())
	}
	if err := pool.Ping(ctx); err != nil {
		panic("banco de teste indisponivel — rode: docker compose up db\n" + err.Error())
	}
	testPool = pool

	code := m.Run()

	pool.Close()
	os.Exit(code)
}

func truncateTables(t *testing.T) {
	t.Helper()
	_, err := testPool.Exec(context.Background(),
		`TRUNCATE TABLE pacientes, usuarios RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatalf("truncate: %v", err)
	}
}