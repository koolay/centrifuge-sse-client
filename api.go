package sseclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type API struct {
	httpclient *http.Client
	url        string
	endpoint   string
}

func NewAPI(url string, options ...func(*API)) *API {
	pub := &API{
		url: url,
	}

	for _, opt := range options {
		opt(pub)
	}
	if pub.httpclient == nil {
		pub.httpclient = http.DefaultClient
	}

	return pub
}

func WithAPIKey(apikey string) func(*API) {
	return func(p *API) {
		p.endpoint = apikey
	}
}

func WithHTTPClient(httpclient *http.Client) func(*API) {
	return func(p *API) {
		p.httpclient = httpclient
	}
}

func (p API) newRequest(params interface{}, apiMethod APIMethod) (*http.Request, error) {
	apiRequest := APIRequest{
		Method: string(apiMethod),
		Params: params,
	}

	body, err := json.Marshal(apiRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, p.url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	if p.endpoint != "" {
		req.Header.Set("Authorization", "apikey "+p.endpoint)
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (p *API) doRequest(ctx context.Context, req *http.Request) ([]byte, error) {
	req = req.WithContext(ctx)
	resp, err := p.httpclient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to http request with url: %s, error: %w", p.url, err)
	}
	defer func() { _ = resp.Body.Close() }()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &ErrBadStatus{Code: resp.StatusCode, Body: resp.Status}
	}

	return respData, nil
}

func (p *API) Publish(ctx context.Context, params *PublishParams) error {
	req, err := p.newRequest(params, publishMethod)
	if err != nil {
		return err
	}

	respData, err := p.doRequest(ctx, req)
	if err != nil {
		return err
	}

	var publishResp PublishResponse
	err = json.Unmarshal(respData, &publishResp)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if publishResp.Error.Code != 0 {
		return &publishResp.Error
	}

	return nil
}

// GetOnlines allows getting channel online presence information (all clients currently subscribed on this channel)
func (p *API) GetOnlines(ctx context.Context, params *GetOnlineParams) (map[string]PresenceItem, error) {
	req, err := p.newRequest(params, presenceMethod)
	if err != nil {
		return nil, err
	}

	respData, err := p.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	var onlineResp GetOnlineResponse
	err = json.Unmarshal(respData, &onlineResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if onlineResp.Error.Code != 0 {
		return nil, &onlineResp.Error
	}

	return onlineResp.Result.Presence, nil
}
