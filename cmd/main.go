package main

import (
	"context"
	"fmt"
	"os"
	"wallet-api/internal/config"
	"wallet-api/internal/database/postgres"
	"wallet-api/internal/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.NewLogger(cfg.Env)
	log.Info("Logger created")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDatabase)

	conn, err := postgres.Connect(connStr, log)
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	log.Info("Server start")
}
