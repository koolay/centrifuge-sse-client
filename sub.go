package sseclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/centrifugal/protocol"
	"github.com/go-logr/logr"
	"github.com/golang-jwt/jwt"

	sse "github.com/r3labs/sse/v2"
	"gopkg.in/cenkalti/backoff.v1"
)

var (
	defaultEventBufferSize = 16
	defaultMaxBufferSize   = 1024 * 1024 * 64
)

type Sub struct {
	logger            *logr.Logger
	sseURL            string
	jwtKey            string
	maxBufferSize     int
	eventBufferSize   int
	reconnectStrategy backoff.BackOff
	OnReconnectNotify backoff.Notify
	OnConnect         sse.ConnCallback
	OnDisconnect      sse.ConnCallback
}

func NewSub(logger *logr.Logger, sseURL string, jwtKey string) *Sub {
	if sseURL == "" {
		sseURL = "http://localhost:8000/connection/uni_sse"
	}

	return &Sub{logger: logger, sseURL: sseURL, jwtKey: jwtKey}
}

func (s *Sub) SetMaxBufferSize(maxBufferSize int) {
	s.maxBufferSize = maxBufferSize
}

func (s *Sub) SetEventBufferSize(eventBufferSize int) {
	s.eventBufferSize = eventBufferSize
}

func (s *Sub) SetReconnectStrategy(reconnectStrategy backoff.BackOff) {
	s.reconnectStrategy = reconnectStrategy
}

func (s *Sub) defaultReconnectStrategy() backoff.BackOff {
	retry := backoff.NewExponentialBackOff()
	retry.MaxInterval = 10 * time.Second
	return retry
}

func (s *Sub) initialize() {
	if s.maxBufferSize == 0 {
		s.maxBufferSize = defaultMaxBufferSize
	}
	if s.eventBufferSize == 0 {
		s.eventBufferSize = defaultEventBufferSize
	}
}

func (s *Sub) prepareClient(channel string, user string) (*sse.Client, error) {
	subReqs := make(map[string]*protocol.SubscribeRequest)
	subReqs[channel] = &protocol.SubscribeRequest{
		// Whether a client wants to recover from a certain position
		Recover: false,
		// Known stream position epoch when recover is used
		Epoch: "",
		// Known stream position offset when recover is used
		Offset: 0,
	}
	token := s.genToken(channel, user, s.jwtKey, 0)
	req := &protocol.ConnectRequest{
		Name:  "1",
		Token: token,
		Subs:  subReqs,
	}
	data, err := json.Marshal(&req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal connect request: %w", err)
	}

	params := url.Values{}
	params.Add("cf_connect", string(data))
	// https://github.com/centrifugal/centrifugo/blob/0be4d975085ac45d1deaba9bca70d091feced2e9/internal/unisse/handler.go#L34
	sseURL := fmt.Sprintf("%s?%s", s.sseURL, params.Encode())
	client := sse.NewClient(sseURL, sse.ClientMaxBufferSize(s.maxBufferSize))
	client.Connection.Transport = &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		TLSHandshakeTimeout: 10 * time.Second,
		Dial: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		}).Dial,
	}

	if s.reconnectStrategy != nil {
		client.ReconnectStrategy = s.reconnectStrategy
	} else {
		client.ReconnectStrategy = s.defaultReconnectStrategy()
	}

	client.ReconnectNotify = func(err error, next time.Duration) {
		s.logger.Info("Reconnecting after", next, "due to", "error", err)
	}
	if s.OnReconnectNotify != nil {
		client.ReconnectNotify = s.OnReconnectNotify
	}

	client.OnConnect(func(c *sse.Client) {
		s.logger.Info("connected to server")
	})
	if s.OnConnect != nil {
		client.OnConnect(s.OnConnect)
	}

	client.OnDisconnect(func(c *sse.Client) {
		log.Println("disconnect from server")
	})
	if s.OnDisconnect != nil {
		client.OnDisconnect(s.OnDisconnect)
	}

	return client, nil
}

func (s *Sub) genToken(channel, user, jwtKey string, exp int64) string {
	// https://centrifugal.dev/docs/transports/uni_sse
	// https://centrifugal.dev/docs/transports/uni_websocket#connect-command
	subs := claimsSub{
		channel: SubscribeOptions{},
	}
	claims := jwt.MapClaims{"sub": user, "subs": subs}
	if exp > 0 {
		claims["exp"] = exp
	}
	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (s *Sub) Subscribe(ctx context.Context, channel string, user string) (<-chan *sse.Event, error) {
	s.initialize()
	client, err := s.prepareClient(channel, user)
	if err != nil {
		return nil, err
	}

	events := make(chan *sse.Event, s.eventBufferSize)
	err = client.SubscribeChanWithContext(ctx, "", events)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w, sseURL: %s", err, s.sseURL)
	}

	return events, nil
}
