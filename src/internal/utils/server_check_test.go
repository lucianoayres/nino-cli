// internal/utils/server_check_test.go
package utils

import (
	"net"
	"testing"
)

func TestIsOllamaRunning(t *testing.T) {
	tests := []struct {
		name     string
		urlStr   string
		want     bool
		setup    func() func()
	}{
		{
			name:   "Valid URL with running server",
			urlStr: "http://localhost:8080",
			want:   true,
			setup: func() func() {
				// Start a TCP server on localhost:8080
				ln, err := net.Listen("tcp", "localhost:8080")
				if err != nil {
					t.Fatalf("Failed to start test server: %v", err)
				}
				// Return a function to close the listener
				return func() {
					ln.Close()
				}
			},
		},
		{
			name:   "Valid URL with no server running",
			urlStr: "http://localhost:8081",
			want:   false,
			setup:  func() func() { return func() {} }, // No setup needed
		},
		{
			name:   "Invalid URL",
			urlStr: "://invalid-url",
			want:   false,
			setup:  func() func() { return func() {} }, // No setup needed
		},
		{
			name:   "URL with empty host",
			urlStr: "http:///path",
			want:   false,
			setup:  func() func() { return func() {} }, // No setup needed
		},
		{
			name:   "Malformed URL",
			urlStr: "http://",
			want:   false,
			setup:  func() func() { return func() {} }, // No setup needed
		},
	}

	for _, tt := range tests {
		// Run each test case in a separate subtest
		t.Run(tt.name, func(t *testing.T) {
			// Setup the environment if needed
			teardown := tt.setup()
			defer teardown()

			// Call the function under test
			got := IsOllamaRunning(tt.urlStr)

			// Check the result
			if got != tt.want {
				t.Errorf("IsOllamaRunning(%q) = %v; want %v", tt.urlStr, got, tt.want)
			}
		})
	}
}
