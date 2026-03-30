package config

import (
	"fmt"
	"os"
	"time"

	"github.com/go-core-fx/config"
)

type database struct {
	URL string `koanf:"url"`

	ConnMaxIdleTime time.Duration `koanf:"conn_max_idle_time"`
	ConnMaxLifetime time.Duration `koanf:"conn_max_lifetime"`
	MaxOpenConns    int           `koanf:"max_open_conns"`
	MaxIdleConns    int           `koanf:"max_idle_conns"`
}

type http struct {
	Address     string   `koanf:"address"`
	ProxyHeader string   `koanf:"proxy_header"`
	Proxies     []string `koanf:"proxies"`

	OpenAPI openAPIConfig `koanf:"openapi"`
}

type openAPIConfig struct {
	Enabled    bool   `koanf:"enabled"`
	PublicHost string `koanf:"public_host"`
	PublicPath string `koanf:"public_path"`
}

type telegram struct {
	Token string `koanf:"token"`
}

type exampleConfig struct {
	Example string `koanf:"example"`
}

type Config struct {
	Database database `koanf:"database"`
	HTTP     http     `koanf:"http"`
	Telegram telegram `koanf:"telegram"`

	Example exampleConfig `koanf:"example"`
}

func Default() Config {
	//nolint:mnd,gosec // default values
	return Config{
		Database: database{
			URL: "mariadb://bot:bot@127.0.0.1:3306/bot?charset=utf8mb4&parseTime=True&loc=UTC",

			ConnMaxIdleTime: 10 * time.Minute,
			ConnMaxLifetime: 1 * time.Hour,
			MaxOpenConns:    25,
			MaxIdleConns:    5,
		},
		HTTP: http{
			Address:     "127.0.0.1:3000",
			ProxyHeader: "X-Forwarded-For",
			Proxies:     []string{},
			OpenAPI: openAPIConfig{
				Enabled:    true,
				PublicHost: "",
				PublicPath: "",
			},
		},

		Telegram: telegram{
			Token: "",
		},

		Example: exampleConfig{
			Example: "example",
		},
	}
}

func New() (Config, error) {
	cfg := Default()

	options := []config.Option{}
	if yamlPath := os.Getenv("CONFIG_PATH"); yamlPath != "" {
		options = append(options, config.WithLocalYAML(yamlPath))
	}

	if err := config.Load(&cfg, options...); err != nil {
		return Config{}, fmt.Errorf("failed to load config: %w", err)
	}

	return cfg, nil
}
