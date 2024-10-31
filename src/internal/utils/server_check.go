package utils

import (
	"net"
	"net/url"
	"time"

	"github.com/lucianoayres/nino-cli/internal/logger"
)

// IsOllamaRunning checks if the Ollama server is running at the specified URL
func IsOllamaRunning(urlStr string) bool {
	log := logger.GetLogger(true) // Assuming logger is already initialized in main
	log.Info("Checking if Ollama server is running at URL: %s", urlStr)

	u, err := url.Parse(urlStr)
	if err != nil {
		log.Error("URL parsing error: %v", err)
		return false
	}
	host := u.Host
	if host == "" {
		log.Error("Empty host in URL")
		return false
	}
	// Try to establish a TCP connection to the host and port
	log.Info("Attempting TCP connection to %s", host)
	conn, err := net.DialTimeout("tcp", host, 2*time.Second)
	if err != nil {
		log.Error("TCP connection failed: %v", err)
		return false
	}
	conn.Close()
	log.Info("TCP connection successful")
	return true
}
