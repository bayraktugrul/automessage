package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type client struct {
	baseURL    string
	httpClient *http.Client
}

type Client interface {
	SendMessage(ctx context.Context, req Request) (messageResponse Response, err error)
}

func New(baseURL string) Client {
	return &client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *client) SendMessage(ctx context.Context, req Request) (messageResponse Response, err error) {
	jsonPayload, err := json.Marshal(req)
	if err != nil {
		return messageResponse, fmt.Errorf("error marshaling request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return messageResponse, fmt.Errorf("error creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return messageResponse, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return messageResponse, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&messageResponse); err != nil {
		return messageResponse, fmt.Errorf("error decoding response: %w", err)
	}

	return messageResponse, nil
}
