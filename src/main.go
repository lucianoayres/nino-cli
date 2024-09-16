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
	// Define the model argument with a default value and the prompt argument with a flag
	modelPtr := flag.String("model", "llama3.1", "The model to use (default is llama3.1)")
	promptPtr := flag.String("prompt", "Explain me LLMs like I'm five", "The prompt to send (required)")

	// Parse the command-line flags
	flag.Parse()

	// Check if the prompt is provided
	if *promptPtr == "" {
		fmt.Println("Usage: query-ollama [-model model_name] -p \"prompt\"")
		os.Exit(1)
	}

	model := *modelPtr
	prompt := *promptPtr
	url := "http://localhost:11434/api/generate"

	// Construct the JSON data with the model and prompt arguments
	data := fmt.Sprintf(`{
		"model": "%s",
		"prompt": "%s"
	}`, model, prompt)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")

	// HTTP client with no timeout
	client := &http.Client{
		Timeout: 0, // No timeout
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	for {
		var r Response
		if err := decoder.Decode(&r); err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error decoding response:", err)
			break
		}

		fmt.Print(r.Response)

		if r.Done {
			break
		}
	}
}
