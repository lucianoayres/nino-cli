package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"nino/internal/client"
	"nino/internal/config"
	"nino/internal/contextmanager"
	"nino/internal/models"
	"nino/internal/processor"
	"nino/internal/utils"
	"os"
	"path/filepath"
)

func main() {
	// Parse command-line arguments using the config package
	cfg, err := config.ParseArgs()
	if err != nil {
		log.Fatalf("Error parsing arguments: %v", err)
	}

	// Check if Ollama server is running
	if !utils.IsOllamaRunning(cfg.URL) {
		fmt.Printf("Oops! It looks like the Ollama server isn't running at %s.\n", cfg.URL)
		fmt.Println("Please start the server at this URL or update the NINO_URL environment variable with the correct URL.")
		fmt.Println("To start the server, you can run:")
		fmt.Printf("ollama serve & ollama run %s\n", cfg.Model)
		os.Exit(1)
	}

	// Initialize the HTTP client
	cli := client.NewHTTPClient(cfg.URL)

	// Read and encode images
	var imagesBase64 []string
	if len(cfg.ImagePaths) > 0 {
		imagesBase64, err = utils.ReadImagesAsBase64(cfg.ImagePaths)
		if err != nil {
			log.Fatalf("Error processing images: %v", err)
		}
	}

	// Prepare the request payload
	payload := models.RequestPayload{
		Model:  cfg.Model,
		Prompt: cfg.Prompt,
		Images: imagesBase64, // Assign the base64-encoded images
		Format: cfg.Format,
		Stream: cfg.Stream,
	}

	// Load context data for the model
	contextData, err := contextmanager.LoadContext(cfg.Model)
	if err != nil {
		log.Fatalf("Error loading context data: %v", err)
	}

	// If context data exists, include it in the payload
	if !cfg.DisableContext && len(contextData) > 0 {
		payload.Context = contextData
	}

	// Start the loading animation in a goroutine if not disabled and not in silent mode
	done := make(chan bool)
	if !cfg.DisableLoading && !cfg.Silent {
		go utils.ShowLoadingAnimation(done)
	}

	// Send the HTTP request
	response, err := cli.SendRequest(payload)

	// Stop the loading animation
	if !cfg.DisableLoading && !cfg.Silent {
		done <- true
	}
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
	var writers []io.Writer
	if !cfg.Silent {
		writers = append(writers, os.Stdout) // Write to console unless in silent mode
	}

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

	// Clear the line before writing the response if not in silent mode
	if !cfg.Silent {
		fmt.Print("\r\033[K")
	}

	// Define context handler
	contextHandler := func(context []int) error {
		return contextmanager.SaveContext(cfg.Model, context)
	}

	// Process the response and write to all writers
	if err := processor.ProcessResponse(response.Body, multiWriter, contextHandler); err != nil {
		log.Fatalf("Error processing response: %v", err)
	}

	// If output was saved to a file and not in silent mode, notify the user
	if cfg.Output != "" && !cfg.Silent {
		fmt.Printf("\nOutput saved to %s\n", cfg.Output)
	} else if !cfg.Silent {
		// Add a newline for console output, so the shell prompt is displayed below
		fmt.Fprintln(os.Stdout)
	}
}
