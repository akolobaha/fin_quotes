package main

import (
	"context"
	"fin_quotes/cmd/commands"
	"fin_quotes/internal/config"
	"fin_quotes/internal/log"
	"fin_quotes/internal/monitoring"
	"fin_quotes/internal/transport"
	"os"
	"os/signal"
	"syscall"
)

const defaultEnvFilePath = "./config.toml"

func init() {
	monitoring.RegisterPrometheus()
}

func main() {
	cfg, err := config.Parse(defaultEnvFilePath)
	if err != nil {
		panic("Ошибка парсинга конфигов")
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		<-exit
		cancel()
	}()

	log.Info("Сервис слежения за котировами запущен")

	rabbit := transport.New()
	rabbit.InitConn(cfg)
	defer rabbit.ConnClose()
	rabbit.DeclareQueue(cfg.RabbitQueue)

	monitoring.RunPrometheusServer(cfg.GetPrometheusURL())

	commands.NewServeCmd(ctx, cfg, rabbit)

}
