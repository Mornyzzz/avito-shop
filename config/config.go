package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		Env  string `yaml:"env"`
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		PG   `yaml:"postgres"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true" yaml:"url"      env:"PG_URL"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	const op = "config.NewConfig"
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
