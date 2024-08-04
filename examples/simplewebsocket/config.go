package main

import "github.com/xdoubleu/essentia/pkg/config"

type Config struct {
	Env            string
	Port           int
	AllowedOrigins []string
}

func NewConfig() Config {
	var cfg Config

	cfg.Env = config.EnvStr("ENV", config.ProdEnv)
	cfg.Port = config.EnvInt("PORT", 8000)
	cfg.AllowedOrigins = config.EnvStrArray(
		"ALLOWED_ORIGINS",
		[]string{"http://localhost"},
	)

	return cfg
}
