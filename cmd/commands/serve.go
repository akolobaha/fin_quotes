package commands

import (
	"context"
	"encoding/json"
	"fin_quotes/internal/config"
	"fin_quotes/internal/quotes"
	"fin_quotes/internal/transport"
	"log/slog"
	"sync"
	"time"
)

func NewServeCmd(ctx context.Context, config *config.Config, rabbit *transport.Rabbitmq) {
	ticker := time.NewTicker(config.Tick)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			data, err := quotes.Fetch(config.Moex)
			if err != nil {
				slog.Error(err.Error())
			} else {
				slog.Info("Данные с Мосбиржи получены")
			}

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

		case <-ctx.Done():
			slog.Info("Сбор данных остановлен")
			return
		}
	}
}
