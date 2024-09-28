// client/client.go
package client

import (
	"bytes"
	"encoding/json"
	"net/http"

	"nino/internal/models"
)

// HTTPClient represents a client for making HTTP requests.
type HTTPClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewHTTPClient creates a new HTTPClient with the given base URL.
func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 0, // No timeout for streaming
		},
	}
}

// SendRequest sends a POST request with the given payload and returns the HTTP response.
func (c *HTTPClient) SendRequest(payload models.RequestPayload) (*http.Response, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
