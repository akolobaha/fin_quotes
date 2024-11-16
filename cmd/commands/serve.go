package commands

import (
	"context"
	"encoding/xml"
	"fin_quotes/internal/config"
	"fin_quotes/internal/quotes"
	"fin_quotes/internal/transport"
	"github.com/spf13/cobra"
	"log/slog"
	"strconv"
	"sync"
	"time"
)

func NewServeCmd(ctx context.Context, config *config.Config, rabbit *transport.Rabbitmq, log *slog.Logger) *cobra.Command {
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
					slog.Info("Данные с Мосбиржи получены")
					if err != nil {
						slog.Error(err.Error())
					}

					//grpc.SendQuotes(ctx, data, config)

					var wg sync.WaitGroup

					for idx, row := range data.Rows {
						wg.Add(1)
						func() {
							defer wg.Done()
							bytes, err := xml.Marshal(row)
							if err != nil {
								slog.Error(err.Error())
							}

							rabbit.SendMsg(bytes)
							slog.Info(strconv.Itoa(idx), string(bytes))
						}()
					}
					wg.Wait()
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
