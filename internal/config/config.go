package config

import (
	"fmt"
	"log"
	"os"
	"wallet-api/internal/logger"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	AmountIsNotValid  = "amount is not valid"
	InvalidOperation  = "invalid operation"
	OperationDeposit  = "DEPOSIT"
	OperationWithdraw = "WITHDRAW"
)

type Config struct {
	Env        string `env:"env" env-default:"local"`
	HTTPServer `env:"http_server"`
	Postgres   `env:"postgres"`
}

type HTTPServer struct {
	Address string `env:"address" env-default:"localhost"`
	Port    string `env:"port" env-default:"80"`
}

type Postgres struct {
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT" env-default:"5432"`
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresDatabase string `env:"POSTGRES_DB"`
	MaxConn          string `env:"postgres_db_pool_max_conns" env-default:"20"`
}

func MustLoad() *Config {
	op := "config.MustLoad"

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("Config path is empty")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("[Error] Cannot read config: ", logger.Err(err), "op", op)
	}

	return &cfg
}

func (cfg *Config) PrintAll() {
	fmt.Println("-------------")
	fmt.Println("Env: " + cfg.Env)
	fmt.Println("Address: " + cfg.Address)
	fmt.Println("Port: " + cfg.Port)
	fmt.Println("Postgres Host: " + cfg.PostgresHost)
	fmt.Println("Postgres Port: " + cfg.PostgresPort)
	fmt.Println("Postgres User: " + cfg.PostgresUser)
	fmt.Println("Postgres Database: " + cfg.PostgresDatabase)
	fmt.Println("-------------")
}
