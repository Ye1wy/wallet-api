package main

import (
	"context"
	"fmt"
	"os"
	"wallet-api/internal/config"
	"wallet-api/internal/controller"
	"wallet-api/internal/logger"
	"wallet-api/internal/repository"
	"wallet-api/internal/routes"
	"wallet-api/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

// @title			Swagger for wallet API
// @version		1.0
// @description	This api created for the test task
// @host			localhost:8080
// @BasePath		/
func main() {
	cfg := config.MustLoad()
	cfg.PrintAll()
	log := logger.NewLogger(cfg.Env)
	log.Info("Logger created")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?pool_max_conns=%s", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDatabase, cfg.MaxConn)
	fmt.Println(connStr)
	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	repos := repository.NewPostgresWalletRepository(dbpool, log)
	log.Info("Repository is created")
	service := service.NewWalletServiceImpl(repos, log)
	service.StartWorkers(5)
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
