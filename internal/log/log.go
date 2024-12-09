package log

import (
	"fin_quotes/internal/monitoring"
	"fmt"
	"log/slog"
)

func Error(additionalMessage string, err error) {
	if err != nil {
		msg := fmt.Sprintf("%s: %s", additionalMessage, err.Error())
		monitoring.QuotesErrorCount.WithLabelValues(msg).Inc()
		slog.Error(msg)
	}
}

func Info(message string) {
	monitoring.QuotesSuccessCount.WithLabelValues(message).Inc()
	slog.Info(message)
}
