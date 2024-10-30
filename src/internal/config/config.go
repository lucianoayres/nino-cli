// config/config.go
package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Config holds the configuration for the request
type Config struct {
	Model          string
	Prompt         string
	PromptFile     string
	URL            string
	Output         string
	DisableLoading bool
	Stream         bool
	Keep_Alive	   string
	DisableContext bool
	Silent         bool
	ImagePaths     []string // New field for image paths
	Format         string
}

// arrayFlags is a custom type for parsing multiple -image flags
type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// ParseArgs parses command-line arguments and returns a Config struct
func ParseArgs() (*Config, error) {
	// Check for environment variables
	defaultModel := os.Getenv("NINO_MODEL")
	if defaultModel == "" {
		defaultModel = "llama3.2" // Fallback default
	}

	defaultURL := os.Getenv("NINO_URL")
	if defaultURL == "" {
		defaultURL = "http://localhost:11434/api/generate" // Fallback default
	}

	defaultKeepAlive := os.Getenv("NINO_KEEP_ALIVE")
	if defaultKeepAlive == "" {
		defaultKeepAlive = "60m" // Keep Model Alive time in minutes
	}

	systemPrompt := os.Getenv("NINO_SYSTEM_PROMPT")

	// Define the flags with their long forms
	modelPtr := flag.String("model", defaultModel, "The model to use (default is llama3.2)")
	promptPtr := flag.String("prompt", "", "The prompt to send (required)")
	promptFilePtr := flag.String("prompt-file", "", "The path to a file containing the prompt (optional)")
	urlPtr := flag.String("url", defaultURL, "The URL to send the request to (default is http://localhost:11434/api/generate)")
	outputPtr := flag.String("output", "", "The file to save the output to (optional)")
	disableLoadingPtr := flag.Bool("no-loading", false, "Disable the loading animation (optional)")
	disableStreamPtr := flag.Bool("no-stream", false, "Disable streaming the output (optional)")
	disableContextPtr := flag.Bool("no-context", false, "Disable the context from the previous request (optional)")
	silentPtr := flag.Bool("silent", false, "Run in silent mode (no console output, requires -output)")
	formatPtr := flag.String("format", "", "The format of the output (must be 'json')")

	// Define short forms for the existing flags
	flag.StringVar(modelPtr, "m", defaultModel, "The model to use (short form)")
	flag.StringVar(promptPtr, "p", "", "The prompt to send (short form, required)")
	flag.StringVar(promptFilePtr, "pf", "", "The file containing the prompt (short form, optional)")
	flag.StringVar(urlPtr, "u", defaultURL, "The URL to send the request to (short form)")
	flag.StringVar(outputPtr, "o", "", "The file to save the output to (short form, optional)")
	flag.BoolVar(disableLoadingPtr, "nl", false, "Disable the loading animation (short form)")
	flag.BoolVar(disableStreamPtr, "ns", false, "Disable streaming the output (short form)")
	flag.BoolVar(disableContextPtr, "nc", false, "Disable the context from the previous request (short form)")
	flag.BoolVar(silentPtr, "s", false, "Run in silent mode (short form, requires -output)")
	flag.StringVar(formatPtr, "f", "", "The format of the output (short form, must be 'json')")

	// Define the new -image flag which can be specified multiple times
	imagePaths := arrayFlags{}

	flag.Var(&imagePaths, "image", "Paths to local image files (can be specified multiple times)")
	flag.Var(&imagePaths, "i", "Paths to local image files (short form)")

	// Customize the usage message (optional)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Parse the flags
	flag.Parse()

	// Validate flags
	if *silentPtr && *outputPtr == "" {
		return nil, errors.New("the -silent flag requires the -output flag to be specified")
	}

	if *formatPtr != "" && *formatPtr != "json" {
		return nil, errors.New("the -format flag must be set to 'json' if specified")
	}

	// If the prompt is not provided via flags, check positional arguments
	if *promptPtr == "" && *promptFilePtr == "" {
		args := flag.Args()
		if len(args) == 0 {
			return nil, errors.New("either the prompt or prompt file is required")
		}
		*promptPtr = strings.Join(args, " ")
	}

	// If the prompt-file is provided, read the file content
	if *promptPtr == "" && *promptFilePtr != "" {
		content, err := os.ReadFile(*promptFilePtr)
		if err != nil {
			return nil, fmt.Errorf("error reading prompt file '%s': %v", *promptFilePtr, err)
		}
		*promptPtr = string(content)
	}

	// Concatenate system prompt if it is set
	if systemPrompt != "" {
		*promptPtr = systemPrompt + " " + *promptPtr
	}

	// Concatenate image paths to the prompt if any
	if len(imagePaths) > 0 {
		*promptPtr = strings.TrimSpace(*promptPtr + " " + strings.Join(imagePaths, " "))
	}

	// Return the Config struct with all fields populated
	return &Config{
		Model:          *modelPtr,
		Prompt:         *promptPtr,
		PromptFile:     *promptFilePtr,
		URL:            *urlPtr,
		Output:         *outputPtr,
		DisableLoading: *disableLoadingPtr,
		Stream:         !*disableStreamPtr,
		Keep_Alive:     defaultKeepAlive,
		DisableContext: *disableContextPtr,
		Silent:         *silentPtr,
		ImagePaths:     imagePaths, // Assign the collected image paths
		Format:         *formatPtr,
	}, nil
}