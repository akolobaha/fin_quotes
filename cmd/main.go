package main

import (
	"context"
	"fin_quotes/cmd/commands"
	"fin_quotes/internal/config"
	pb "fin_quotes/pkg/grpc"

	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

const defaultConfigPath = "../config.yaml"

var Log *slog.Logger

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	fmt.Println(cfg.Tick)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		<-exit
		cancel()
	}()

	log.Info("Сервис запущен")

	// TODO: пробросить данные с апихи мосбиржи
	grpcFn()

	cmd := commands.NewServeCmd(cfg, ctx, log)

	cmd.ExecuteContext(ctx)
}

func grpcFn() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewDataManagementServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	currentTime := time.Now()
	timestamp := timestamppb.New(currentTime)

	msg := &pb.TickerRequest{Name: "hello", Price: "hello", Time: timestamp}

	response, err := c.GetQuotes(ctx, msg)
	if err != nil {
		log.Fatalf("could not send message: %v", err)
	}

	log.Printf("Response from server: %s", response.Time, response.Price, response.Name)

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
