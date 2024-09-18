// config/config.go
package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// Config holds the configuration for the request
type Config struct {
	Model  string
	Prompt string
	URL    string
	Output string // Field to hold the output file path
}

// ParseArgs parses command-line arguments and returns a Config struct
func ParseArgs() (*Config, error) {
	// Define the flags with their long forms
	modelPtr := flag.String("model", "llama3.1", "The model to use (default is llama3.1)")
	promptPtr := flag.String("prompt", "", "The prompt to send (required)")
	urlPtr := flag.String("url", "http://localhost:11434/api/generate", "The URL to send the request to (default is http://localhost:11434/api/generate)")
	outputPtr := flag.String("output", "", "The file to save the output to (optional)")

	// Define short forms for the existing flags
	flag.StringVar(modelPtr, "m", "llama3.1", "The model to use (short form)")
	flag.StringVar(promptPtr, "p", "", "The prompt to send (short form, required)")
	flag.StringVar(urlPtr, "u", "http://localhost:11434/api/generate", "The URL to send the request to (short form)")
	flag.StringVar(outputPtr, "o", "", "The file to save the output to (short form, optional)")

	// Customize the usage message (optional)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Parse the flags
	flag.Parse()

	// Check if the required prompt is provided
	if *promptPtr == "" {
		return nil, errors.New("the prompt is required")
	}

	// Return the Config struct with all fields populated
	return &Config{
		Model:  *modelPtr,
		Prompt: *promptPtr,
		URL:    *urlPtr,
		Output: *outputPtr, // Assign the output file path
	}, nil
}
