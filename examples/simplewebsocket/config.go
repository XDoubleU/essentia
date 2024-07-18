package main

import "github.com/XDoubleU/essentia/pkg/config"

type Config struct {
	Env            string
	Port           int
	AllowedOrigins []string
}

const (
	ProdEnv = "production"
	TestEnv = "test"
)

func NewConfig() Config {
	var cfg Config

	cfg.Env = config.GetEnvStr("ENV", ProdEnv)
	cfg.Port = config.GetEnvInt("PORT", 8000)
	cfg.AllowedOrigins = config.GetEnvStrArray(
		"ALLOWED_ORIGINS",
		[]string{"http://localhost"},
	)

	return cfg
}
