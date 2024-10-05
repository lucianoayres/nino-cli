package config

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantConfig *Config
		wantErr    bool
	}{
		{
			name: "Valid arguments with long flags",
			args: []string{"cmd", "--model=llama3.1", "--prompt=Hello", "--url=http://localhost:11434/api/generate", "--output=result.txt"},
			wantConfig: &Config{
				Model:  "llama3.1",
				Prompt: "Hello",
				URL:    "http://localhost:11434/api/generate",
				Output: "result.txt",
			},
			wantErr: false,
		},
		{
			name: "Valid arguments with short flags",
			args: []string{"cmd", "-m", "llama3.1", "-p", "Hello", "-u", "http://localhost:11434/api/generate", "-o", "result.txt"},
			wantConfig: &Config{
				Model:  "llama3.1",
				Prompt: "Hello",
				URL:    "http://localhost:11434/api/generate",
				Output: "result.txt",
			},
			wantErr: false,
		},
		{
			name:    "Missing required prompt argument",
			args:    []string{"cmd", "--model=llama3.1"},
			wantErr: true,
		},
		{
			name: "Default values for optional arguments",
			args: []string{"cmd", "--prompt=Hello"},
			wantConfig: &Config{
				Model:  "llama3.1",
				Prompt: "Hello",
				URL:    "http://localhost:11434/api/generate",
				Output: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save the original os.Args and restore it after the test
			origArgs := os.Args
			defer func() {
				os.Args = origArgs
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) // Reset flag.CommandLine to avoid flag redefinition errors
			}()

			// Set os.Args to the test case arguments
			os.Args = tt.args

			// Parse the arguments
			gotConfig, err := ParseArgs()
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(gotConfig, tt.wantConfig) {
				t.Errorf("ParseArgs() = %v, want %v", gotConfig, tt.wantConfig)
			}
		})
	}
}