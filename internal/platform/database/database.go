package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	defaultMaxConns     = 10
	defaultConnAttempts = 3
	defaultConnTimeout  = 5 * time.Second
)

type Config struct {
	URL          string
	MaxConns     int32
	ConnAttempts int
	ConnTimeout  time.Duration
}

func NewConfig() *Config {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "sidekick"
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	return &Config{
		URL:          url,
		MaxConns:     defaultMaxConns,
		ConnAttempts: defaultConnAttempts,
		ConnTimeout:  defaultConnTimeout,
	}
}

func Connect(ctx context.Context) (*pgxpool.Pool, error) {
	cfg := NewConfig()

	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parsing database URL: %w", err)
	}

	poolConfig.MaxConns = cfg.MaxConns

	var pool *pgxpool.Pool
	for i := 0; i < cfg.ConnAttempts; i++ {
		connCtx, cancel := context.WithTimeout(ctx, cfg.ConnTimeout)
		pool, err = pgxpool.ConnectConfig(connCtx, poolConfig)
		cancel()
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	return pool, nil
}
