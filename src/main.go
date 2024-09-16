package main

import (
	"bytes"
	"encoding/json"
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
	url := "http://localhost:11434/api/generate"
	data := `{
		"model": "llama3.1",
		"prompt": "Why is the sky blue? I want you to give a really short answer!"
	}`

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
