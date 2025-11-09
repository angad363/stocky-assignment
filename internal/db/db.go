package db

import (
	"fmt"

	"github.com/angad363/stocky-assignment/internal/config"
	"github.com/angad363/stocky-assignment/pkg/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(cfg *config.Config) *sqlx.DB {
	dbURL := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	logger.Log.Infof("Connecting to PostgreSQL at %s:%s...", cfg.DBHost, cfg.DBPort)
	conn, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		logger.Log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	logger.Log.Info("Connected to PostgreSQL successfully")
	return conn
}
