package integration

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestMigrationsCreated_UsersTableExists(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("connect error: %v", err)
	}
	defer pool.Close()

	var exists bool
	err = pool.QueryRow(ctx, "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name='users')").Scan(&exists)
	if err != nil {
		t.Fatalf("query error: %v", err)
	}
	if !exists {
		t.Fatalf("expected users table to exist after migrations")
	}
}
