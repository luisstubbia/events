package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func NewDBClient(ctx context.Context, cfg *Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.GetConnectionString())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func getConnectionString(config *Config) string {
	// Use the traditional connection string format that works with lib/pq
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
		config.SSLMode,
	)
}
