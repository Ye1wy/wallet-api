package config

type Config struct {
	Env string `env:"env" env-default:"local"`
}
