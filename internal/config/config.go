package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	Env         string        `toml:"env" env-default:"local"`
	Tick        time.Duration `toml:"tick" env-default:"1"`
	Moex        string        `toml:"moex"`
	DataService string        `toml:"dataService"`

	SourceUrl      string `toml:"SOURCE_URL"`
	RabbitUser     string `toml:"RABBIT_USER"`
	RabbitPassword string `toml:"RABBIT_PASSWORD"`
	RabbitHost     string `toml:"RABBIT_HOST"`
	RabbitPort     int64  `toml:"RABBIT_PORT"`
	RabbitQueue    string `toml:"RABBIT_QUEUE"`
}

func Parse(s string) (*Config, error) {
	c := &Config{}
	if err := cleanenv.ReadConfig(s, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) GetRabbitDSN() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d/", c.RabbitUser, c.RabbitPassword, c.RabbitHost, c.RabbitPort,
	)
}
