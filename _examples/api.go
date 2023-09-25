package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-logr/logr"
	sseclient "github.com/koolay/centrifuge-sse-client"
)

type Post struct {
	Number int `json:"number,omitempty"`
}

// https://centrifugal.dev/docs/3/server/channels#presence
// should enable presence, else this won't work('not available')
func getOnline(logger *logr.Logger) {
	p := sseclient.NewAPI("http://localhost:8000/api", sseclient.WithAPIKey(apiKey))
	result, err := p.GetOnlines(context.Background(), &sseclient.GetOnlineParams{
		Channel: channelTest,
	})
	if err != nil {
		logger.Error(err, "faied to get online")
	}

	logger.Info("got onlines", "result", result)
}

func publish(ctx context.Context, logger *logr.Logger) {
	var i int
	for {
		i++
		data, _ := json.Marshal(Post{Number: i})
		p := sseclient.NewAPI("http://localhost:8000/api", sseclient.WithAPIKey(apiKey))
		err := p.Publish(context.Background(), &sseclient.PublishParams{
			Channel:     channelTest,
			Data:        data,
			SkipHistory: true,
		})
		if err != nil {
			logger.Error(err, "failed to publish")
		}

		time.Sleep(time.Second)
	}
}
