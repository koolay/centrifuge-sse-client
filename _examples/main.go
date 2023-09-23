package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"

	"github.com/koolay/centrifuge-sse-client/sseclient"
)

var (
	channelTest = "facts:devops"

	userid = "1"

	apiKey = os.Getenv("SSE_API_KEY")
	jwtKey = os.Getenv("SSE_JWT_KEY")
)

func getLogger() logr.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerologr.SetMaxV(1)
	zl := zerolog.New(os.Stdout)
	return zerologr.New(&zl)
}

func main() {
	sseURL := "http://localhost:8000/connection/uni_sse"
	logger := getLogger()
	sseClient := sseclient.NewSub(&logger, sseURL, jwtKey)
	events, err := sseClient.Subscribe(context.Background(), channelTest, userid)
	if err != nil {
		panic(err)
	}

	go func() {
		for event := range events {
			logger.Info("event", "data", event.Data, "event", event.Event, "offset", event.ID)
		}
	}()

	go publish(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	<-sigChan
}
