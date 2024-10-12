package processor

import (
	"bytes"
	"errors"
	"testing"
)

// Helper function to compare two slices of integers
func slicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// TestProcessResponse tests the ProcessResponse function with various input scenarios.
func TestProcessResponse(t *testing.T) {
	tests := []struct {
		name                string
		input               string
		wantOutput          string
		wantErr             bool
		expectedContext     []int
		contextHandlerError error
	}{
		{
			name: "Normal processing with context",
			input: `{"Response": "Hello", "Done": false}
{"Response": " World", "Done": true, "Context": [1,2,3]}`,
			wantOutput:      "Hello World",
			wantErr:         false,
			expectedContext: []int{1, 2, 3},
		},
		{
			name: "Skip empty responses with context",
			input: `{"Response": "", "Done": false}
{"Response": "Valid Response", "Done": true, "Context": [4,5,6]}`,
			wantOutput:      "Valid Response",
			wantErr:         false,
			expectedContext: []int{4, 5, 6},
		},
		{
			name: "Malformed JSON input",
			input: `{"Response": "Hello", "Done": false}
{"Response": "World", "Done": true`,
			wantOutput: "",
			wantErr:    true,
		},
		{
			name:            "Empty input",
			input:           ``,
			wantOutput:      "",
			wantErr:         false,
			expectedContext: nil,
		},
		{
			name: "Only empty responses without context",
			input: `{"Response": "", "Done": false}
{"Response": "", "Done": true}`,
			wantOutput:      "",
			wantErr:         false,
			expectedContext: nil,
		},
		{
			name: "Single response marked as done with context",
			input: `{"Response": "Single Response", "Done": true, "Context": [42]}`,
			wantOutput:      "Single Response",
			wantErr:         false,
			expectedContext: []int{42},
		},
		{
			name: "Multiple responses with early done",
			input: `{"Response": "First", "Done": false}
{"Response": "Second", "Done": true, "Context": [7,8]}
{"Response": "Third", "Done": false}`,
			wantOutput:      "FirstSecond",
			wantErr:         false,
			expectedContext: []int{7, 8},
		},
		{
			name: "Response with special characters",
			input: `{"Response": "Hello\nWorld!", "Done": true, "Context": [9]}`,
			wantOutput:      "Hello\nWorld!",
			wantErr:         false,
			expectedContext: []int{9},
		},
		{
			name: "Response with unicode characters",
			input: `{"Response": "こんにちは", "Done": true, "Context": [10]}`,
			wantOutput:      "こんにちは",
			wantErr:         false,
			expectedContext: []int{10},
		},
		{
			name: "Context handler returns error",
			input: `{"Response": "Test", "Done": true, "Context": [99]}`,
			wantOutput:          "Test",
			wantErr:             true,
			expectedContext:     []int{99},
			contextHandlerError: errors.New("context handler error"),
		},
		{
			name: "No context provided",
			input: `{"Response": "No Context", "Done": true}`,
			wantOutput:      "No Context",
			wantErr:         false,
			expectedContext: nil,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create an io.Reader from the input string
			reader := bytes.NewReader([]byte(tt.input))

			// Use a bytes.Buffer as the io.Writer to capture the output
			var writer bytes.Buffer

			// Variables to capture context handler calls
			var receivedContext []int
			var contextHandlerCalled bool

			// Define the contextHandler function
			contextHandler := func(context []int) error {
				contextHandlerCalled = true
				receivedContext = context
				return tt.contextHandlerError
			}

			// Call the ProcessResponse function
			err := ProcessResponse(reader, &writer, contextHandler)

			// Check if an error was expected
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If an error was expected from the context handler, verify it
			if tt.contextHandlerError != nil && err != nil {
				expectedErrMsg := "failed to handle context: " + tt.contextHandlerError.Error()
				if err.Error() != expectedErrMsg {
					t.Errorf("ProcessResponse() error = %v, expected %v", err.Error(), expectedErrMsg)
				}
			}

			// If no error was expected, compare the output
			if !tt.wantErr {
				gotOutput := writer.String()
				if gotOutput != tt.wantOutput {
					t.Errorf("ProcessResponse() output = %q, want %q", gotOutput, tt.wantOutput)
				}
			}

			// Verify if contextHandler was called as expected
			if tt.expectedContext != nil {
				if !contextHandlerCalled {
					t.Errorf("Expected contextHandler to be called, but it was not")
				} else if !slicesEqual(receivedContext, tt.expectedContext) {
					t.Errorf("Context received = %v, expected %v", receivedContext, tt.expectedContext)
				}
			} else if contextHandlerCalled {
				t.Errorf("Did not expect contextHandler to be called, but it was")
			}
		})
	}
}
