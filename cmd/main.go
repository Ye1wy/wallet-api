package main

import (
	"context"
	"fmt"
	"os"
	"wallet-api/internal/config"
	"wallet-api/internal/controller"
	"wallet-api/internal/database/postgres"
	"wallet-api/internal/logger"
	"wallet-api/internal/repository"
	"wallet-api/internal/routes"
	"wallet-api/internal/service"
)

func main() {
	cfg := config.MustLoad()
	cfg.PrintAll()
	log := logger.NewLogger(cfg.Env)
	log.Info("Logger created")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDatabase)
	fmt.Println(connStr)
	conn, err := postgres.Connect(connStr, log)
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	repos := repository.NewPostgresWalletRepository(conn, log)
	log.Info("Repository is created")
	service := service.NewWalletServiceImpl(repos, log)
	log.Info("Service is created")
	controller := controller.NewWalletController(service, log)
	log.Info("Controller is created")

	routerConfig := routes.RouterConfig{
		WalletController: controller,
	}

	router := routes.NewRouter(routerConfig)
	log.Info("Router is created")
	log.Info("Server start")
	router.Run(":" + cfg.Port)
}
