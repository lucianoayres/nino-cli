package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"gollama/internal/client"
	"gollama/internal/config"
	"gollama/internal/models"
	"gollama/internal/processor"
)

func main() {
    cfg, err := config.ParseArgs()
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }

    cli := client.NewHTTPClient(cfg.URL)
    payload := models.RequestPayload{
        Model:  cfg.Model,
        Prompt: cfg.Prompt,
    }

    response, err := cli.SendRequest(payload)
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        bodyBytes, _ := io.ReadAll(response.Body)
        fmt.Printf("Error: Received HTTP status %d\nResponse body: %s\n", response.StatusCode, string(bodyBytes))
        os.Exit(1)
    }

    if err := processor.ProcessResponse(response.Body); err != nil {
        fmt.Println("Error processing response:", err)
        os.Exit(1)
    }
}
