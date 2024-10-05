package utils

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestShowLoadingAnimation(t *testing.T) {
	// Save the original stdout
	originalStdout := os.Stdout
	defer func() {
		os.Stdout = originalStdout
	}()

	// Create a pipe to capture output
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	// Redirect stdout to the pipe
	os.Stdout = w

	// Create a channel to signal completion
	done := make(chan bool)

	// Start the loading animation in a goroutine
	go ShowLoadingAnimation(done)

	// Let the animation run for a short time
	time.Sleep(1000 * time.Millisecond)

	// Signal the loading animation to stop
	done <- true

	// Wait a moment to ensure the goroutine has exited
	time.Sleep(100 * time.Millisecond)

	// Close the write end of the pipe and restore stdout
	w.Close()

	// Read the captured output
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}

	output := buf.String()

	// Remove ANSI escape codes from the output
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	cleanOutput := re.ReplaceAllString(output, "")

	// Check that the output is not empty
	if len(cleanOutput) == 0 {
		t.Error("Expected output, but got empty string")
	}

	// Check that the output contains the loading text
	if !strings.Contains(cleanOutput, "Thinking") {
		t.Errorf("Expected output to contain 'Thinking', but got: %s", cleanOutput)
	}
}
