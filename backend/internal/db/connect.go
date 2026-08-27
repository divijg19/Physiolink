package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/divijg19/physiolink/backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB is the runtime DB wrapper used by the application. It holds both a
// pgxpool for code that expects pgx and a *sql.DB for sqlc-generated code.
type DB struct {
	Pool    *pgxpool.Pool
	SQL     *sql.DB
	Queries *Queries
}

// Connect establishes connections (pgxpool and database/sql via pgx stdlib)
// and returns a DB wrapper ready for use.
func Connect(ctx context.Context, cfg *config.Config) (*DB, error) {
	if cfg == nil {
		return nil, fmt.Errorf("nil config")
	}

	// create pgxpool for existing code that uses pgx
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}

	// create database/sql DB using pgx stdlib driver (registered by import)
	sqlDB, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// reasonable defaults
	sqlDB.SetConnMaxLifetime(60 * time.Minute)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)

	// ping to verify connectivity
	if err := sqlDB.PingContext(ctx); err != nil {
		pool.Close()
		_ = sqlDB.Close()
		return nil, fmt.Errorf("sqlDB.Ping: %w", err)
	}

	// create sqlc queries wrapper
	queries := New(sqlDB)

	return &DB{Pool: pool, SQL: sqlDB, Queries: queries}, nil
}

// Close releases DB resources.
func (d *DB) Close() error {
	if d == nil {
		return nil
	}

	if d.Pool != nil {
		d.Pool.Close()
	}
	if d.SQL != nil {
		return d.SQL.Close()
	}
	return nil
}
