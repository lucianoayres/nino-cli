// File: internal/client/client_test.go

package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gollama/internal/models"
)

func TestHTTPClient_SendRequest_Success(t *testing.T) {
	// Create a test server that responds with a predefined response
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate HTTP method
		if r.Method != http.MethodPost {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// Validate Content-Type header
		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", contentType)
		}

		// Read and validate request body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		defer r.Body.Close()

		expectedBody := `{"model":"test-model","prompt":"test-prompt"}`
		if strings.TrimSpace(string(bodyBytes)) != expectedBody {
			t.Errorf("Expected request body %s, got %s", expectedBody, string(bodyBytes))
		}

		// Send a successful JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response":"test-response","done":true}`))
	}))
	defer testServer.Close()

	// Create an instance of HTTPClient with the test server URL
	client := NewHTTPClient(testServer.URL)

	// Prepare the request payload
	payload := models.RequestPayload{
		Model:  "test-model",
		Prompt: "test-prompt",
	}

	// Send the request
	resp, err := client.SendRequest(payload)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	// Validate response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	// Read and validate response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	expectedResponse := `{"response":"test-response","done":true}`
	if strings.TrimSpace(string(respBody)) != expectedResponse {
		t.Errorf("Expected response body %s, got %s", expectedResponse, string(respBody))
	}
}

func TestHTTPClient_SendRequest_NonOKStatus(t *testing.T) {
	// Create a test server that responds with a 500 Internal Server Error
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer testServer.Close()

	client := NewHTTPClient(testServer.URL)
	payload := models.RequestPayload{
		Model:  "test-model",
		Prompt: "test-prompt",
	}

	resp, err := client.SendRequest(payload)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	// Validate response status code
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code 500, got %d", resp.StatusCode)
	}

	// Read and validate response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	expectedResponse := "Internal Server Error"
	if strings.TrimSpace(string(respBody)) != expectedResponse {
		t.Errorf("Expected response body %s, got %s", expectedResponse, string(respBody))
	}
}

func TestHTTPClient_SendRequest_RequestCreationError(t *testing.T) {
	// Use an invalid URL to trigger request creation error
	client := NewHTTPClient("://invalid-url")

	payload := models.RequestPayload{
		Model:  "test-model",
		Prompt: "test-prompt",
	}

	_, err := client.SendRequest(payload)
	if err == nil {
		t.Fatalf("Expected error, got none")
	}

	if !strings.Contains(err.Error(), "failed to create request") {
		t.Errorf("Expected request creation error, got %v", err)
	}
}

func TestHTTPClient_SendRequest_DoRequestError(t *testing.T) {
	// Use an unreachable URL to trigger an error during Do request
	client := NewHTTPClient("http://localhost:0") // Port 0 is invalid

	payload := models.RequestPayload{
		Model:  "test-model",
		Prompt: "test-prompt",
	}

	_, err := client.SendRequest(payload)
	if err == nil {
		t.Fatalf("Expected error, got none")
	}

	if !strings.Contains(err.Error(), "failed to send request") {
		t.Errorf("Expected request sending error, got %v", err)
	}
}
