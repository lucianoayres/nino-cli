// main.go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gollama/internal/client"
	"gollama/internal/config"
	"gollama/internal/models"
	"gollama/internal/processor"
)

func main() {
	// Parse command-line arguments
	cfg, err := config.ParseArgs()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Initialize the HTTP client
	cli := client.NewHTTPClient(cfg.URL)

	// Prepare the request payload
	payload := models.RequestPayload{
		Model:  cfg.Model,
		Prompt: cfg.Prompt,
	}

	// Send the HTTP request
	response, err := cli.SendRequest(payload)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer response.Body.Close()

	// Check for non-OK HTTP status
	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(response.Body)
		log.Fatalf("Error: Received HTTP status %d\nResponse body: %s", response.StatusCode, string(bodyBytes))
	}

	// Prepare writers
	writers := []io.Writer{os.Stdout} // Always write to console

	// If Output is specified, add the file to writers
	if cfg.Output != "" {
		// Validate the output directory exists
		dir := filepath.Dir(cfg.Output)
		if dir != "." { // Skip if current directory
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				log.Fatalf("Error: Directory '%s' does not exist.", dir)
			}
		}

		file, err := os.Create(cfg.Output)
		if err != nil {
			log.Fatalf("Error creating output file '%s': %v", cfg.Output, err)
		}
		defer file.Close()
		writers = append(writers, file)
	}

	// Create a MultiWriter to write to all destinations
	multiWriter := io.MultiWriter(writers...)

	// Process the response and write to all writers
	if err := processor.ProcessResponse(response.Body, multiWriter); err != nil {
		log.Fatalf("Error processing response: %v", err)
	}

	// If output was saved to a file, notify the user
	if cfg.Output != "" {
		fmt.Printf("\nOutput saved to %s\n", cfg.Output)
	}
}
