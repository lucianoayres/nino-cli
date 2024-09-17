package config

import (
	"errors"
	"flag"
)

// Config holds the configuration for the request
type Config struct {
    Model  string
    Prompt string
    URL    string
}

// ParseArgs parses command-line arguments and returns a Config struct
func ParseArgs() (*Config, error) {
    modelPtr := flag.String("model", "llama3.1", "The model to use (default is llama3.1)")
    promptPtr := flag.String("prompt", "", "The prompt to send (required)")
    urlPtr := flag.String("url", "http://localhost:11434/api/generate", "The URL to send the request to (default is http://localhost:11434/api/generate)")

    flag.StringVar(modelPtr, "m", "llama3.1", "The model to use (short form)")
    flag.StringVar(promptPtr, "p", "", "The prompt to send (short form, required)")
    flag.StringVar(urlPtr, "u", "http://localhost:11434/api/generate", "The URL to send the request to (short form)")

    flag.Parse()

    if *promptPtr == "" {
        return nil, errors.New("the prompt is required")
    }

    return &Config{
        Model:  *modelPtr,
        Prompt: *promptPtr,
        URL:    *urlPtr,
    }, nil
}
