package commands

import (
	"context"
	"encoding/json"
	"fin_quotes/internal/config"
	"fin_quotes/internal/quotes"
	"fin_quotes/internal/transport"
	"github.com/spf13/cobra"
	"log/slog"
	"sync"
	"time"
)

func NewServeCmd(ctx context.Context, config *config.Config, rabbit *transport.Rabbitmq) *cobra.Command {
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
					data, err := quotes.Fetch(config.Moex)
					if err != nil {
						//slog.Error(err)
					} else {
						slog.Info("Данные с Мосбиржи получены")
					}
					//
					////grpc.SendQuotes(ctx, data, config)
					//
					var wg sync.WaitGroup

					for _, row := range data {
						wg.Add(1)
						func() {
							defer wg.Done()
							bytes, err := json.Marshal(row)
							if err != nil {
								slog.Error(err.Error())
							}

							rabbit.SendMsg(bytes)
						}()
					}
					slog.Info("Котировки отправленны в брокер сообщений")

					//wg.Wait()
				case <-ctx.Done():
					slog.Info("Сбор данных остановлен")
					return nil
				}
			}
		},
	}

	c.Flags().StringVar(&configPath, "config", "", "path to config")
	return c
}
