package main

import (
	"wallet-api/internal/config"
	"wallet-api/internal/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.NewLogger(cfg.Env)

	log.Info("Server start")
}
