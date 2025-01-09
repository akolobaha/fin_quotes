package commands

import (
	"context"
	"encoding/json"
	"fin_quotes/internal/config"
	"fin_quotes/internal/log"
	"fin_quotes/internal/quotes"
	"fin_quotes/internal/transport"
	"fmt"
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
				log.Error("Ошибка получения данных с API мосбиржи", err)
			} else {
				log.Info("Данные с Мосбиржи получены")
			}

			var wg sync.WaitGroup

			for _, row := range data {
				wg.Add(1)
				func() {
					defer wg.Done()
					bytes, err := json.Marshal(row)
					if err != nil {
						log.Error("", err)
					}

					rabbit.SendMsg(bytes)
				}()
			}
			log.Info(fmt.Sprintf("Котировки отправленны в брокер сообщений %d шт", len(data)))

		case <-ctx.Done():
			log.Info("Сбор данных остановлен")
			return
		}
	}
}
