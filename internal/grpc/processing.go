package grpc

import (
	"context"
	"fin_quotes/internal/config"
	"fin_quotes/internal/quotes"
	pb "fin_quotes/pkg/grpc"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

func SendQuotes(ctx context.Context, marketData quotes.MarketData, cfg *config.Config) {
	conn, err := grpc.Dial(cfg.DataService, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewDataManagementServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reqItems := make([]*pb.TickerRequest, len(marketData.Rows))
	for _, item := range marketData.Rows {

		t, err := time.Parse("15:04:05", item.TIME)
		if err != nil {
			log.Fatalf("Error parsing time: %v", err)
		}
		timestamp := timestamppb.New(t)

		reqItem := pb.TickerRequest{Name: item.SECID, Price: item.LAST, Time: timestamp}
		reqItems = append(reqItems, &reqItem)
	}

	msg := &pb.MultipleTickerRequest{Tickers: reqItems}

	response, err := c.GetMultipleQuotes(ctx, msg)
	if err != nil {
		log.Fatalf("could not send message: %v", err)
	}

	fmt.Println("Response from server: %s", response)
	//log.Printf("Response from server: %s", response.Time, response.Price, response.Name)

}
