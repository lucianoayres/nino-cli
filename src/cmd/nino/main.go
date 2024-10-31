package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/lucianoayres/nino-cli/internal/client"
	"github.com/lucianoayres/nino-cli/internal/config"
	"github.com/lucianoayres/nino-cli/internal/contextmanager"
	"github.com/lucianoayres/nino-cli/internal/logger"
	"github.com/lucianoayres/nino-cli/internal/models"
	"github.com/lucianoayres/nino-cli/internal/processor"
	"github.com/lucianoayres/nino-cli/internal/utils"
)

func main() {
	// Parse command-line arguments using the config package
	cfg, err := config.ParseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing arguments: %v\n", err)
		os.Exit(1)
	}

	// Initialize the logger
	log := logger.GetLogger(cfg.Verbose)

	log.StartTimer("Total Execution Time")
	defer log.StopTimer("Total Execution Time")

	log.Info("Starting NINO CLI tool")

	// Check if Ollama server is running
	log.StartTimer("Check Ollama Server")
	log.Info("Checking if Ollama server is running at %s", cfg.URL)
	if !utils.IsOllamaRunning(cfg.URL) {
		fmt.Printf("Oops! It looks like the Ollama server isn't running at %s.\n", cfg.URL)
		fmt.Println("Please start the server at this URL or update the NINO_URL environment variable with the correct URL.")
		fmt.Println("To start the server, you can run:")
		fmt.Printf("ollama serve & ollama run %s\n", cfg.Model)
		os.Exit(1)
	}
	log.Info("Ollama server is running")
	log.StopTimer("Check Ollama Server")

	// Initialize the HTTP client
	log.StartTimer("Initialize HTTP Client")
	log.Info("Initializing HTTP client with base URL: %s", cfg.URL)
	cli := client.NewHTTPClient(cfg.URL)
	log.StopTimer("Initialize HTTP Client")

	// Read and encode images
	var imagesBase64 []string
	if len(cfg.ImagePaths) > 0 {
		log.StartTimer("Process Images")
		log.Info("Reading and encoding %d image(s)", len(cfg.ImagePaths))
		imagesBase64, err = utils.ReadImagesAsBase64(cfg.ImagePaths)
		if err != nil {
			log.Error("Error processing images: %v", err)
			os.Exit(1)
		}
		log.Info("Images processed successfully")
		log.StopTimer("Process Images")
	} else {
		log.Info("No images provided")
	}

	// Prepare the request payload
	log.StartTimer("Prepare Request Payload")
	log.Info("Preparing request payload")
	payload := models.RequestPayload{
		Model:      cfg.Model,
		Prompt:     cfg.Prompt,
		Images:     imagesBase64, // Assign the base64-encoded images
		Format:     cfg.Format,
		Stream:     cfg.Stream,
		Keep_Alive: cfg.Keep_Alive,
	}
	log.StopTimer("Prepare Request Payload")

	// TODO: Fix the performance for context data 
	// Justification: It's slowing down the application performance

	/*
	// Load context data for the model
	contextData, err := contextmanager.LoadContext(cfg.Model)
	if err != nil {
		log.Fatalf("Error loading context data: %v", err)
	}

	
	// If context data exists, include it in the payload
	if !cfg.DisableContext && len(contextData) > 0 {
		payload.Context = contextData
	}
	*/
	
	// Start the loading animation in a goroutine if not disabled and not in silent mode
	done := make(chan bool)
	if !cfg.DisableLoading && !cfg.Silent {
		go utils.ShowLoadingAnimation(done)
	}

	// Send the HTTP request
	log.StartTimer("Send HTTP Request")
	log.Info("Sending HTTP request to Ollama server")
	response, err := cli.SendRequest(payload)
	log.StopTimer("Send HTTP Request")

	// Stop the loading animation
	if !cfg.DisableLoading && !cfg.Silent {
		done <- true
	}

	if err != nil {
		log.Error("Error sending request: %v", err)
		os.Exit(1)
	}
	defer response.Body.Close()
	log.Info("Received response with status code: %d", response.StatusCode)

	// Check for non-OK HTTP status
	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(response.Body)
		log.Error("Error: Received HTTP status %d\nResponse body: %s", response.StatusCode, string(bodyBytes))
		os.Exit(1)
	}
	log.Info("HTTP request successful")

	// Prepare writers
	var writers []io.Writer
	if !cfg.Silent {
		writers = append(writers, os.Stdout) // Write to console unless in silent mode
	}

	// If Output is specified, add the file to writers
	if cfg.Output != "" {
		log.StartTimer("Prepare Output File")
		log.Info("Output will be saved to file: %s", cfg.Output)
		// Validate the output directory exists
		dir := filepath.Dir(cfg.Output)
		if dir != "." { // Skip if current directory
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				log.Error("Error: Directory '%s' does not exist.", dir)
				os.Exit(1)
			}
		}

		file, err := os.Create(cfg.Output)
		if err != nil {
			log.Error("Error creating output file '%s': %v", cfg.Output, err)
			os.Exit(1)
		}
		defer file.Close()
		writers = append(writers, file)
		log.Info("Output file created successfully")
		log.StopTimer("Prepare Output File")
	}

	// Create a MultiWriter to write to all destinations
	multiWriter := io.MultiWriter(writers...)

	// Clear the line before writing the response if not in silent mode
	if !cfg.Silent {
		fmt.Print("\r\033[K")
	}

	// Define context handler
	contextHandler := func(context []int) error {
		log.StartTimer("Save Context Data")
		log.Info("Saving context data")
		err := contextmanager.SaveContext(cfg.Model, context)
		if err != nil {
			log.Error("Failed to save context data: %v", err)
			log.StopTimer("Save Context Data")
			return err
		}
		log.StopTimer("Save Context Data")
		return nil
	}

	// Process the response and write to all writers
	log.StartTimer("Process Response")
	log.Info("Processing response")
	if err := processor.ProcessResponse(response.Body, multiWriter, contextHandler); err != nil {
		log.Error("Error processing response: %v", err)
		os.Exit(1)
	}
	log.Info("Response processed successfully")
	log.StopTimer("Process Response")

	// If output was saved to a file and not in silent mode, notify the user
	if cfg.Output != "" && !cfg.Silent {
		fmt.Printf("\nOutput saved to %s\n", cfg.Output)
		log.Info("Output saved to file")
	} else if !cfg.Silent {
		// Add a newline for console output, so the shell prompt is displayed below
		fmt.Fprintln(os.Stdout)
	}
	log.Info("NINO CLI tool completed successfully")
}
