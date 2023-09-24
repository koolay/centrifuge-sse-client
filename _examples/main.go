package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"time"

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
	sseClient := sseclient.NewSub(&logger, sseURL,
		sseclient.WithJWTKey(jwtKey),
	)
	events, err := sseClient.Subscribe(context.Background(), channelTest, userid)
	if err != nil {
		panic(err)
	}

	go func() {
		for event := range events {
			var eventData sseclient.EventData
			err := json.Unmarshal(event.Data, &eventData)
			if err != nil {
				logger.Error(err, "failed to unmarshal event data")
				continue
			}

			logger.Info("received", "data", eventData, "event", event.Event, "eventID", event.ID)
		}
	}()

	go publish(context.Background(), &logger)
	time.Sleep(time.Second)
	go getOnline(&logger)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	<-sigChan
}
