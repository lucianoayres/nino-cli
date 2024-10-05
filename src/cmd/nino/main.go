package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"nino/internal/client"
	"nino/internal/models"
	"nino/internal/processor"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// Check for the environment variable "NINO_MODEL"
	defaultModel := os.Getenv("NINO_MODEL")
	if defaultModel == "" {
		defaultModel = "llama3.2" // Fallback default if the environment variable is not set
	}

	// Check for the environment variable "NINO_URL"
	defaultURL := os.Getenv("NINO_URL")
	if defaultURL == "" {
		defaultURL = "http://localhost:11434/api/generate" // Fallback default if the environment variable is not set
	}

	// Define command-line flags
	model := flag.String("model", defaultModel, "The model to use")
	prompt := flag.String("prompt", "", "The prompt to send to the language model")
	promptFile := flag.String("prompt-file", "", "The path to a file containing the prompt to send to the language model")
	url := flag.String("url", defaultURL, "The host and port where the Ollama server is running")
	output := flag.String("output", "", "Specifies the filename where the model output will be saved")
	disableLoading := flag.Bool("no-loading", false, "Disable the loading animation")

	// Parse the flags
	flag.Parse()

	// Check if the prompt is provided via -prompt or -prompt-file
	if *prompt == "" && *promptFile == "" {
		args := flag.Args()
		if len(args) == 0 {
			fmt.Println("Error: No prompt provided. Use -prompt, -prompt-file flag or provide prompt as positional arguments.")
			os.Exit(1)
		}
		*prompt = strings.Join(args, " ")
	}

	// If the prompt-file is provided, read the file content
	if *prompt == "" && *promptFile != "" {
		content, err := os.ReadFile(*promptFile)
		if err != nil {
			log.Fatalf("Error reading prompt file '%s': %v", *promptFile, err)
		}
		*prompt = string(content)
	}

	// **Check if Ollama server is running**
	if !isOllamaRunning(*url) {
		fmt.Println("Oops! It looks like the Ollama server isn't running.")
		fmt.Println("Please start the server and run the model you'd like to use. For example:")
		fmt.Printf("ollama serve & ollama run %s\n", *model)
		os.Exit(1)
	}

	// Initialize the HTTP client
	cli := client.NewHTTPClient(*url)

	// Prepare the request payload
	payload := models.RequestPayload{
		Model:  *model,
		Prompt: *prompt,
	}

	// Start the loading animation in a goroutine if not disabled
	done := make(chan bool)
	if !*disableLoading {
		go showLoadingAnimation(done)
	}

	// Send the HTTP request
	response, err := cli.SendRequest(payload)

	// Stop the loading animation
	if !*disableLoading {
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

	// Clear the line before writing the response
	fmt.Print("\r\033[K")

	// Process the response and write to all writers
	if err := processor.ProcessResponse(response.Body, multiWriter); err != nil {
		log.Fatalf("Error processing response: %v", err)
	}

	// If output was saved to a file, notify the user
	if *output != "" {
		fmt.Printf("\nOutput saved to %s\n", *output)
	} else {
		// Add a newline for console output, so the shell prompt is displayed below
		fmt.Fprintln(os.Stdout)
	}
}

// **Function to check if Ollama server is running**
func isOllamaRunning(urlStr string) bool {
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

// showLoadingAnimation displays a loading animation in the console
func showLoadingAnimation(done chan bool) {
	words := []string{"Thinking"} // Add more words to the list to pick it randomly
	loadingText := words[rand.Intn(len(words))]
	shades := []string{
		"\033[1;30m", // Dark Gray
		"\033[1;90m", // Light Dark Gray
		"\033[1;37m", // Light Gray
		"\033[0;37m", // White
	}
	resetColor := "\033[0m"
	
	for {
		select {
		case <-done:
			// Clear the animation before stopping
			fmt.Print("\r\033[K")
			return
		default:
			// Create a wave effect by iterating over each character and applying shades
			for waveStart := 0; waveStart < len(loadingText)+len(shades); waveStart++ {
				fmt.Printf("\r")
				for i := 0; i < len(loadingText); i++ {
					shadeOffset := waveStart - i
					if shadeOffset >= 0 && shadeOffset < len(shades) {
						fmt.Printf("%s%c%s", shades[len(shades)-1-shadeOffset], loadingText[i], resetColor)
					} else {
						fmt.Printf("%s%c%s", shades[0], loadingText[i], resetColor)
					}
				}
				time.Sleep(150 * time.Millisecond)
			}
		}
	}
}