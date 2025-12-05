package config

import (
	"os"
)

type Config struct {
	BindAddr    string
	DatabaseURL string
	RedisURL    string
	Env         string
	JWTSecret   string
}

func New() *Config {
	bind := os.Getenv("BIND_ADDR")
	if bind == "" {
		bind = ":8080"
	}
	db := os.Getenv("DATABASE_URL")
	if db == "" {
		db = "postgres://postgres:postgres@localhost:5432/physiolink?sslmode=disable"
	}
	redis := os.Getenv("REDIS_URL")
	if redis == "" {
		redis = "redis://localhost:6379"
	}
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	jwt := os.Getenv("JWT_SECRET")
	if jwt == "" {
		jwt = "changeme"
	}

	return &Config{
		BindAddr:    bind,
		DatabaseURL: db,
		RedisURL:    redis,
		Env:         env,
		JWTSecret:   jwt,
	}
}
