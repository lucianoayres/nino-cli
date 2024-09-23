// main.go
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gollama/internal/client"
	"gollama/internal/models"
	"gollama/internal/processor"
)

func main() {
	// Define command-line flags
	model := flag.String("model", "llama3.1", "The model to use")
	prompt := flag.String("prompt", "", "The prompt to send to the language model")
	url := flag.String("url", "http://localhost:11434/api/generate", "The host and port where the Ollama server is running")
	output := flag.String("output", "", "Specifies the filename where the model output will be saved")
	flag.Parse()

	// If prompt is empty, check for positional arguments
	if *prompt == "" {
		args := flag.Args()
		if len(args) == 0 {
			fmt.Println("Error: No prompt provided. Use -prompt flag or provide prompt as positional arguments.")
			os.Exit(1)
		}
		*prompt = strings.Join(args, " ")
	}

	// Initialize the HTTP client
	cli := client.NewHTTPClient(*url)

	// Prepare the request payload
	payload := models.RequestPayload{
		Model:  *model,
		Prompt: *prompt,
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
	if *output != "" {
		// Validate the output directory exists
		dir := filepath.Dir(*output)
		if dir != "." { // Skip if current directory
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				log.Fatalf("Error: Directory '%s' does not exist.", dir)
			}
		}

		file, err := os.Create(*output)
		if err != nil {
			log.Fatalf("Error creating output file '%s': %v", *output, err)
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
	if *output != "" {
		fmt.Printf("\nOutput saved to %s\n", *output)
	}
}
