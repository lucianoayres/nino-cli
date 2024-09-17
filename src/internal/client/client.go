// File: internal/client/client.go

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"gollama/internal/models"
)

// HTTPClient implements the Client interface using the net/http package
type HTTPClient struct {
    httpClient *http.Client
    url        string
}

// NewHTTPClient creates a new instance of HTTPClient
func NewHTTPClient(url string) *HTTPClient {
    return &HTTPClient{
        httpClient: &http.Client{Timeout: 0}, // No timeout
        url:        url,
    }
}

// SendRequest sends an HTTP POST request with the given payload
func (c *HTTPClient) SendRequest(payload models.RequestPayload) (*http.Response, error) {
    data, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request payload: %v", err)
    }

    req, err := http.NewRequest("POST", c.url, bytes.NewBuffer(data))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %v", err)
    }

    return resp, nil
}
