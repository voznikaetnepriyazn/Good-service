package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv"
)

type Config struct {
	Env        string     `env:"APP_ENV" env-default:"local" env-required:"true"`
	DB         DBConfig   `env-prefix:"DB_"`
	HTTPServer HttpServer `env-prefix:"HTTP_"`
}

type DBConfig struct {
	Host     string `env:"HOST" env-required:"true"`
	Port     int    `env:"PORT" env-default:"5432"`
	User     string `env:"USER" env-default:"postgres"`
	Password string `env:"PASSWORD" env-required:"true"`
	Name     string `env:"NAME" env-required:"true"`
	SSLMode  string `env:"SSLMode" env-default:"disable"`
}

func (db DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.Name, db.SSLMode,
	)
}

type HttpServer struct {
	Address     string `env:"ADDRESS" env-default:"localhost:8080"`
	Timeout     int64  `env:"TIMEOUT" env-default:"4000000000"`
	IdleTimeout int64  `env:"IDLE_TIMEOUT" env-default:"6000000000"`
}

func (h HttpServer) AsDuration() time.Duration {
	return time.Duration(h.Timeout)
}

func (h HttpServer) AsIdleDuration() time.Duration {
	return time.Duration(h.IdleTimeout)
}

func MustLoad() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		slog.Error("failed to load config from env")
	}

	if cfg.Env == "" {
		slog.Error("APP_ENV cannot be empty")
	}
	if cfg.DB.DSN() == "" {
		slog.Error("database configuration is incomplete")
	}

	slog.Info(fmt.Sprintf("config loaded: env=%s, db=%s, http=%s",
		cfg.Env, cfg.DB.Host, cfg.HTTPServer.Address))

	return &cfg
}
