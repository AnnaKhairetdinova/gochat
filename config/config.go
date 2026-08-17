package config

import (
	"github.com/AnnaKhairetdinova/gochat/pkg"
)

type Config struct {
	DNS            string
	Port           string
	HistoryLimit   int
	SendBufferSize int
	RedisAddr      string
	InstanceID     string
}

func Load() Config {
	cfg := Config{}

	cfg.DNS = pkg.MustEnv("DATABASE_URL")
	cfg.Port = pkg.EnvOr("HTTP_PORT", ":8080")
	cfg.HistoryLimit = pkg.EnvInt("HISTORY_LIMIT", 50)
	cfg.SendBufferSize = pkg.EnvInt("SEND_BUFFER_SIZE", 256)
	cfg.RedisAddr = pkg.MustEnv("REDIS_ADDR")
	cfg.InstanceID = pkg.EnvOr("INSTANCE_ID", "unknown")

	return cfg
}
