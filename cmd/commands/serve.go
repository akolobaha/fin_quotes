package commands

import (
	"context"
	"fin_quotes/internal/config"
	"fin_quotes/internal/quotes"
	"github.com/spf13/cobra"
	"log/slog"
	"time"
)

func NewServeCmd(config *config.Config, ctx context.Context, log *slog.Logger) *cobra.Command {
	var configPath string
	c := &cobra.Command{
		Use:     "period",
		Aliases: []string{"s"},
		Short:   "Tick rate period",
		RunE: func(cmd *cobra.Command, args []string) error {
			tick := time.Tick(config.Tick)

			for {
				select {
				case <-tick:
					quotes.Fetch(config.Moex)
					log.Info("Данные с Мосбиржи получены")
				case <-ctx.Done():
					log.Info("Сбор данных остановлен")
					return nil
				}
			}

		},
	}

	c.Flags().StringVar(&configPath, "config", "", "path to config")
	return c
}
