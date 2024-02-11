package env

import (
	"errors"
	"github.com/alexptdev/chat-server-api/internal/config"
	"os"
)

const (
	pgDsnName = "PG_DSN"
)

type pgConfig struct {
	dsn string
}

func NewPgConfig() (config.PgConfig, error) {

	dsn := os.Getenv(pgDsnName)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	return &pgConfig{dsn: dsn}, nil
}

func (cfg pgConfig) Dsn() string {
	return cfg.dsn
}
