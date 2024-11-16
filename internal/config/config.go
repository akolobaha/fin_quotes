package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

const defaultConfigPath = "config.toml"

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
	RabbitDsn      string
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config path is empty: " + err.Error())
	}
	cfg.RabbitDsn = InitRabbitDSN(&cfg)

	return &cfg
}
func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = defaultConfigPath
	}

	return res
}

func InitRabbitDSN(c *Config) string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d/", c.RabbitUser, c.RabbitPassword, c.RabbitHost, c.RabbitPort,
	)
}
