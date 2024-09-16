package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Response struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

func main() {
	// Define flags for model, prompt, and url with default values and short options
	modelPtr := flag.String("model", "llama3.1", "The model to use (default is llama3.1)")
	promptPtr := flag.String("prompt", "", "The prompt to send (required)")
	urlPtr := flag.String("url", "http://localhost:11434/api/generate", "The URL to send the request to (default is http://localhost:11434/api/generate)")

	// Support short forms (-m for --model, -p for --prompt, -u for --url)
	flag.StringVar(modelPtr, "m", "llama3.1", "The model to use (short form)")
	flag.StringVar(promptPtr, "p", "", "The prompt to send (short form, required)")
	flag.StringVar(urlPtr, "u", "http://localhost:11434/api/generate", "The URL to send the request to (short form)")

	// Parse the command-line flags
	flag.Parse()

	// Ensure the prompt is provided, otherwise exit with a usage message
	if *promptPtr == "" {
		fmt.Println("Error: The prompt is required.")
		flag.Usage()
		os.Exit(1)
	}

	model := *modelPtr
	prompt := *promptPtr
	url := *urlPtr

	// Construct the JSON data with the model and prompt arguments
	data := fmt.Sprintf(`{
		"model": "%s",
		"prompt": "%s"
	}`, model, prompt)

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client with no timeout
	client := &http.Client{
		Timeout: 0, // No timeout
	}

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Decode the response
	decoder := json.NewDecoder(resp.Body)

	// Loop to read and print the response until done
	for {
		var r Response
		if err := decoder.Decode(&r); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error decoding response:", err)
			break
		}

		// Print the response
		fmt.Print(r.Response)

		// Break if the response is marked as done
		if r.Done {
			break
		}
	}
}
