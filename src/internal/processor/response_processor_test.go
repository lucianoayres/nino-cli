package processor

import (
	"bytes"
	"testing"
)

// TestProcessResponse tests the ProcessResponse function with various input scenarios.
func TestProcessResponse(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantOutput string
		wantErr    bool
	}{
		{
			name: "Normal processing with multiple responses",
			input: `{"Response": "Hello", "Done": false}
{"Response": " World", "Done": true}`,
			wantOutput: "Hello World",
			wantErr:    false,
		},
		{
			name: "Skip empty responses",
			input: `{"Response": "", "Done": false}
{"Response": "Valid Response", "Done": true}`,
			wantOutput: "Valid Response",
			wantErr:    false,
		},
		{
			name: "Malformed JSON input",
			input: `{"Response": "Hello", "Done": false}
{"Response": "World", "Done": true`,
			wantOutput: "",
			wantErr:    true,
		},
		{
			name:       "Empty input",
			input:      ``,
			wantOutput: "",
			wantErr:    false,
		},
		{
			name: "Only empty responses",
			input: `{"Response": "", "Done": false}
{"Response": "", "Done": true}`,
			wantOutput: "",
			wantErr:    false,
		},
		{
			name: "Single response marked as done",
			input: `{"Response": "Single Response", "Done": true}`,
			wantOutput: "Single Response",
			wantErr:    false,
		},
		{
			name: "Multiple responses with early done",
			input: `{"Response": "First", "Done": false}
{"Response": "Second", "Done": true}
{"Response": "Third", "Done": false}`,
			wantOutput: "FirstSecond",
			wantErr:    false,
		},
		{
			name: "Response with special characters",
			input: `{"Response": "Hello\nWorld!", "Done": true}`,
			wantOutput: "Hello\nWorld!",
			wantErr:    false,
		},
		{
			name: "Response with unicode characters",
			input: `{"Response": "こんにちは", "Done": true}`,
			wantOutput: "こんにちは",
			wantErr:    false,
		},
		{
			name: "Response with nested JSON objects",
			input: `{"Response": "{\"key\":\"value\"}", "Done": true}`,
			wantOutput: "{\"key\":\"value\"}",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create an io.Reader from the input string
			reader := bytes.NewReader([]byte(tt.input))

			// Use a bytes.Buffer as the io.Writer to capture the output
			var writer bytes.Buffer

			// Call the ProcessResponse function
			err := ProcessResponse(reader, &writer)

			// Check if an error was expected
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If no error was expected, compare the output
			if !tt.wantErr {
				gotOutput := writer.String()
				if gotOutput != tt.wantOutput {
					t.Errorf("ProcessResponse() output = %q, want %q", gotOutput, tt.wantOutput)
				}
			}
		})
	}
}
