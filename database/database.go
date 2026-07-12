package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDBPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("Error parsing database URL: %v", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnIdleTime = 10 * time.Minute
	config.MaxConnLifetime = 30 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	dbpool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Error trying to connect to the database: %v", err)
	}

	if err := dbpool.Ping(ctx); err != nil {
		dbpool.Close()
		return nil, fmt.Errorf("Error trying to ping the database: %v", err)
	}

	return dbpool, nil
}
