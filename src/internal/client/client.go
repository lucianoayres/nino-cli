package client

import (
	"bytes"
	"encoding/json"

	"net/http"

	"github.com/lucianoayres/nino-cli/internal/logger"
	"github.com/lucianoayres/nino-cli/internal/models"
)

// HTTPClient represents a client for making HTTP requests.
type HTTPClient struct {
	BaseURL    string
	HTTPClient *http.Client
	log        *logger.Logger
}

// NewHTTPClient creates a new HTTPClient with the given base URL.
func NewHTTPClient(baseURL string) *HTTPClient {
	log := logger.GetLogger(true) // Assuming logger is already initialized in main
	log.Info("Creating new HTTPClient with BaseURL: %s", baseURL)
	return &HTTPClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 0, // No timeout for streaming
		},
		log: log,
	}
}

// SendRequest sends a POST request with the given payload and returns the HTTP response.
func (c *HTTPClient) SendRequest(payload models.RequestPayload) (*http.Response, error) {
	c.log.Info("Marshaling request payload to JSON")
	jsonData, err := json.Marshal(payload)
	if err != nil {
		c.log.Error("JSON marshaling error: %v", err)
		return nil, err
	}
	c.log.Info("JSON payload marshaled successfully")
	
	// Log the request payload
	c.log.Info("Request payload: %s", string(jsonData))

	c.log.Info("Creating new HTTP POST request to %s", c.BaseURL)
	req, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		c.log.Error("HTTP request creation error: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	c.log.Info("HTTP request headers set: Content-Type=application/json")

	// Send the request
	c.log.Info("Sending HTTP request")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		c.log.Error("HTTP request error: %v", err)
		return nil, err
	}

	c.log.Info("Received HTTP response with status code: %d", resp.StatusCode)
	return resp, nil
}
