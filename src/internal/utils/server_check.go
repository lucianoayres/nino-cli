// internal/utils/server_check.go
package utils

import (
	"net"
	"net/url"
	"time"
)

// IsOllamaRunning checks if the Ollama server is running at the specified URL
func IsOllamaRunning(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	host := u.Host
	if host == "" {
		return false
	}
	// Try to establish a TCP connection to the host and port
	conn, err := net.DialTimeout("tcp", host, 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
