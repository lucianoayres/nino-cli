package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"testing"

	"gollama/internal/models"
)

// mockRoundTripper is a custom RoundTripper for mocking HTTP responses.
type mockRoundTripper struct {
	mockResponse *http.Response
	mockError    error
}

// RoundTrip implements the RoundTripper interface.
func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.mockResponse, m.mockError
}

// TestHTTPClient_SendRequest tests the SendRequest method of the HTTPClient.
func TestHTTPClient_SendRequest(t *testing.T) {
	// Define a sample RequestPayload for testing
	samplePayload := models.RequestPayload{
		// Populate with appropriate fields based on your actual RequestPayload struct
		// Example:
		// Name: "TestUser",
		// Age:  25,
	}

	// Helper function to create a mocked HTTP client
	newMockClient := func(rt http.RoundTripper) *HTTPClient {
		return &HTTPClient{
			BaseURL:    "http://mocked-url.com",
			HTTPClient: &http.Client{Transport: rt},
		}
	}

	tests := []struct {
		name           string
		payload        models.RequestPayload
		mockResponse   *http.Response
		mockError      error
		expectedStatus int
		wantErr        bool
	}{
		{
			name:    "Successful POST request",
			payload: samplePayload,
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"Response": "Success", "Done": true}`)),
				Header:     make(http.Header),
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:    "JSON marshaling error",
			payload: models.RequestPayload{
				// Introduce a field that cannot be marshaled, e.g., a channel
				// Uncomment and modify based on your actual RequestPayload struct
				// UnmarshalableField: make(chan int),
			},
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: 0,
			wantErr:        true,
		},
		{
			name:           "HTTP request creation error",
			payload:        samplePayload,
			mockResponse:   nil,
			mockError:      errors.New("invalid URL"),
			expectedStatus: 0,
			wantErr:        true,
		},
		{
			name:    "HTTP client Do request error",
			payload: samplePayload,
			mockResponse: &http.Response{
				// The response won't be used since an error is mocked
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(``)),
				Header:     make(http.Header),
			},
			mockError:      errors.New("network error"),
			expectedStatus: 0,
			wantErr:        true,
		},
		{
			name:    "Server returns non-200 status",
			payload: samplePayload,
			mockResponse: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewBufferString(`{"error": "Internal Server Error"}`)),
				Header:     make(http.Header),
			},
			mockError:      nil,
			expectedStatus: http.StatusInternalServerError,
			wantErr:        false,
		},
		{
			name:    "Response with special characters",
			payload: samplePayload,
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"Response": "Hello\nWorld!", "Done": true}`)),
				Header:     make(http.Header),
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:    "Response with unicode characters",
			payload: samplePayload,
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"Response": "こんにちは", "Done": true}`)),
				Header:     make(http.Header),
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			wantErr:        false,
		},
		{
			name:    "Response with nested JSON objects",
			payload: samplePayload,
			mockResponse: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"Response": "{\"key\":\"value\"}", "Done": true}`)),
				Header:     make(http.Header),
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			wantErr:        false,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Run tests in parallel

			var client *HTTPClient

			// Handle the "HTTP request creation error" test case
			if tt.name == "HTTP request creation error" {
				// Simulate invalid URL by setting a malformed BaseURL
				client = &HTTPClient{
					BaseURL:    "http://[::1]:NamedPort", // Invalid URL to trigger error
					HTTPClient: &http.Client{},
				}
			} else if tt.name == "JSON marshaling error" {
				// To simulate JSON marshaling error, we'll inject a mocked HTTP client
				// that doesn't matter since the error occurs before making the request
				client = newMockClient(&mockRoundTripper{
					mockResponse: nil,
					mockError:    nil,
				})
			} else {
				// For other test cases, use the mocked RoundTripper
				rt := &mockRoundTripper{
					mockResponse: tt.mockResponse,
					mockError:    tt.mockError,
				}
				client = newMockClient(rt)
			}

			// Call the SendRequest method
			resp, err := client.SendRequest(tt.payload)

			// Check if an error was expected
			if (err != nil) != tt.wantErr {
				t.Errorf("SendRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If no error was expected, verify the response
			if !tt.wantErr {
				if tt.expectedStatus != 0 && resp.StatusCode != tt.expectedStatus {
					t.Errorf("Expected status code %d, got %d", tt.expectedStatus, resp.StatusCode)
				}

				// Optionally, read and verify the response body
				if resp != nil && resp.Body != nil {
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						t.Errorf("Failed to read response body: %v", err)
					}

					// For demonstration, unmarshal into ResponsePayload and verify
					var respPayload models.ResponsePayload
					if err := json.Unmarshal(body, &respPayload); err != nil {
						t.Errorf("Failed to unmarshal response body: %v", err)
					}

				}
			}
		})
	}

	// Additional Tests that do not fit into the table-driven approach
	t.Run("Multiple concurrent SendRequest calls", func(t *testing.T) {
		t.Parallel() // Run in parallel with other tests

		// Setup a single mocked RoundTripper for all concurrent requests
		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"Response": "OK", "Done": true}`)),
			Header:     make(http.Header),
		}

		rt := &mockRoundTripper{
			mockResponse: mockResp,
			mockError:    nil,
		}
		client := newMockClient(rt)

		numRequests := 10
		errCh := make(chan error, numRequests)

		for i := 0; i < numRequests; i++ {
			go func() {
				payload := samplePayload
				resp, err := client.SendRequest(payload)
				if err != nil {
					errCh <- err
					return
				}
				defer resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					errCh <- errors.New("unexpected status code")
					return
				}
				errCh <- nil
			}()
		}

		for i := 0; i < numRequests; i++ {
			if err := <-errCh; err != nil {
				t.Errorf("Concurrent SendRequest failed: %v", err)
			}
		}
	})
}

// TestNewHTTPClient tests the NewHTTPClient constructor.
func TestNewHTTPClient(t *testing.T) {
	t.Parallel() // Run in parallel with other tests

	baseURL := "http://example.com/api"
	client := NewHTTPClient(baseURL)

	if client.BaseURL != baseURL {
		t.Errorf("Expected BaseURL %s, got %s", baseURL, client.BaseURL)
	}

	if client.HTTPClient == nil {
		t.Error("Expected HTTPClient to be initialized, got nil")
	}

	// Verify that the Timeout is set to 0
	if client.HTTPClient.Timeout != 0 {
		t.Errorf("Expected HTTPClient Timeout 0, got %v", client.HTTPClient.Timeout)
	}
}

// TestHTTPClient_SendRequest_InvalidURL tests SendRequest with an invalid URL.
func TestHTTPClient_SendRequest_InvalidURL(t *testing.T) {
	t.Parallel() // Run in parallel with other tests

	client := NewHTTPClient("http://[::1]:NamedPort") // Invalid URL

	payload := models.RequestPayload{
		// Populate with appropriate fields
	}

	resp, err := client.SendRequest(payload)
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
	if resp != nil {
		t.Errorf("Expected response to be nil, got %v", resp)
	}
}

// TestHTTPClient_SendRequest_NetworkError tests SendRequest when a network error occurs.
func TestHTTPClient_SendRequest_NetworkError(t *testing.T) {
	t.Parallel() // Run in parallel with other tests

	// Simulate a network error by using a RoundTripper that always returns a network error
	client := &HTTPClient{
		BaseURL: "http://mocked-url.com",
		HTTPClient: &http.Client{
			Transport: &mockRoundTripper{
				mockResponse: nil,
				mockError:    &net.OpError{Op: "dial", Net: "tcp", Addr: nil, Err: errors.New("simulated network error")},
			},
		},
	}

	payload := models.RequestPayload{
		// Populate with appropriate fields
	}

	resp, err := client.SendRequest(payload)
	if err == nil {
		t.Error("Expected network error, got nil")
	}

	// Corrected errors.As usage
	var opErr *net.OpError
	if !errors.As(err, &opErr) {
		t.Logf("Received error: %v", err)
	}

	if resp != nil {
		t.Errorf("Expected response to be nil, got %v", resp)
	}
}
