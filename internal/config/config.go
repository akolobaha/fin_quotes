package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log/slog"
	"time"
)

type Config struct {
	Env         string        `toml:"env" env-default:"local"`
	Tick        time.Duration `toml:"tick" env-default:"1"`
	Moex        string        `toml:"moex"`
	DataService string        `toml:"dataService"`

	SourceURL      string `toml:"SOURCE_URL"`
	RabbitUser     string `toml:"RABBIT_USER"`
	RabbitPassword string `toml:"RABBIT_PASSWORD"`
	RabbitHost     string `toml:"RABBIT_HOST"`
	RabbitPort     int64  `toml:"RABBIT_PORT"`
	RabbitQueue    string `toml:"RABBIT_QUEUE"`
	LogLevel       string `toml:"LOG_LEVEL"`

	PrometheusPort int64  `toml:"PROMETHEUS_PORT"`
	PrometheusHost string `toml:"PROMETHEUS_HOST"`
}

func Parse(s string) (*Config, error) {
	c := &Config{}
	if err := cleanenv.ReadConfig(s, c); err != nil {
		return nil, err
	}
	setLogLevel(c.LogLevel)
	return c, nil
}

func setLogLevel(level string) {
	switch level {
	case "debug":
		slog.SetLogLoggerLevel(-4)
	case "info":
		slog.SetLogLoggerLevel(0)
	case "warn":
		slog.SetLogLoggerLevel(4)
	case "error":
		slog.SetLogLoggerLevel(8)
	default:
		slog.SetLogLoggerLevel(4)
	}
}

func (c *Config) GetRabbitDSN() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d/", c.RabbitUser, c.RabbitPassword, c.RabbitHost, c.RabbitPort,
	)
}

func (c *Config) GetPrometheusURL() string {
	return fmt.Sprintf(
		"%s:%d", c.PrometheusHost, c.PrometheusPort,
	)
}
