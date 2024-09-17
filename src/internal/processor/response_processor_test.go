package processor

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// captureOutput captures the standard output during the execution of a function
func captureOutput(f func()) string {
	var buf bytes.Buffer
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	done := make(chan struct{})
	go func() {
		io.Copy(&buf, r)
		close(done)
	}()

	f()

	w.Close()
	os.Stdout = origStdout
	<-done

	return buf.String()
}

func TestProcessResponse(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput string
		expectedError  string
	}{
		{
			name:           "Empty input",
			input:          ``,
			expectedOutput: ``,
			expectedError:  "",
		},
		{
			name:           "Single response, done true",
			input:          `{"response":"Hello, world!","done":true}`,
			expectedOutput: "Hello, world!",
			expectedError:  "",
		},
		{
			name: "Multiple responses, done false then true",
			input: `{"response":"Hello, ","done":false}
{"response":"world!","done":true}`,
			expectedOutput: "Hello, world!",
			expectedError:  "",
		},
		{
			name: "Invalid JSON",
			input: `{"response":"Hello","done":false}
{"response":,"done":false}`,
			expectedOutput: "Hello",
			expectedError:  "failed to decode JSON response: invalid character ',' looking for beginning of value",
		},
		{
			name: "No done true, ends with EOF",
			input: `{"response":"Hello, ","done":false}
{"response":"world!","done":false}`,
			expectedOutput: "Hello, world!",
			expectedError:  "",
		},
		{
			name: "Error after done true",
			input: `{"response":"Hello, ","done":false}
{"response":"world!","done":true}
{"response":"Extra","done":false}`,
			expectedOutput: "Hello, world!",
			expectedError:  "",
		},
		{
			name: "JSON decoding error",
			input: `{"response":"Hello, ","done":false}
{"response":,"done":false}
{"response":"world!","done":true}`,
			expectedOutput: "Hello, ",
			expectedError:  "failed to decode JSON response: invalid character ',' looking for beginning of value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)

			var actualErr error
			output := captureOutput(func() {
				actualErr = ProcessResponse(reader)
			})

			// Trim spaces for consistent comparison
			output = strings.TrimSpace(output)
			expectedOutput := strings.TrimSpace(tt.expectedOutput)

			// Check for expected error
			if tt.expectedError != "" {
				if actualErr == nil || actualErr.Error() != tt.expectedError {
					t.Errorf("expected error %q, got %v", tt.expectedError, actualErr)
				}
			} else {
				if actualErr != nil {
					t.Errorf("unexpected error: %v", actualErr)
				}
			}

			// Check for expected output
			if output != expectedOutput {
				t.Errorf("expected output %q, got %q", expectedOutput, output)
			}
		})
	}
}
