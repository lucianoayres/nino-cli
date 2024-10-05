// config/config_test.go
package config

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	// Create a temporary prompt file for testing
	tmpDir := t.TempDir()
	promptFilePath := filepath.Join(tmpDir, "prompt.txt")
	err := ioutil.WriteFile(promptFilePath, []byte("Hello from file"), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary prompt file: %v", err)
	}

	tests := []struct {
		name           string
		args           []string
		envModel       string
		envURL         string
		wantConfig     *Config
		wantErr        bool
		wantErrMessage string
	}{
		{
			name: "Valid arguments with long flags",
			args: []string{"cmd", "--model=llama3.1", "--prompt=Hello", "--url=http://localhost:11434/api/generate", "--output=result.txt"},
			wantConfig: &Config{
				Model:          "llama3.1",
				Prompt:         "Hello",
				PromptFile:     "",
				URL:            "http://localhost:11434/api/generate",
				Output:         "result.txt",
				DisableLoading: false,
				Silent:         false,
			},
			wantErr: false,
		},
		{
			name: "Valid arguments with short flags",
			args: []string{"cmd", "-m", "llama3.1", "-p", "Hello", "-u", "http://localhost:11434/api/generate", "-o", "result.txt"},
			wantConfig: &Config{
				Model:          "llama3.1",
				Prompt:         "Hello",
				PromptFile:     "",
				URL:            "http://localhost:11434/api/generate",
				Output:         "result.txt",
				DisableLoading: false,
				Silent:         false,
			},
			wantErr: false,
		},
		{
			name:           "Missing required prompt argument",
			args:           []string{"cmd", "--model=llama3.1"},
			wantErr:        true,
			wantErrMessage: "either the prompt or prompt file is required",
		},
		{
			name: "Default values for optional arguments",
			args: []string{"cmd", "--prompt=Hello"},
			wantConfig: &Config{
				Model:          "llama3.2",
				Prompt:         "Hello",
				PromptFile:     "",
				URL:            "http://localhost:11434/api/generate",
				Output:         "",
				DisableLoading: false,
				Silent:         false,
			},
			wantErr: false,
		},
		{
			name: "Prompt from positional arguments",
			args: []string{"cmd", "Hello", "world"},
			wantConfig: &Config{
				Model:          "llama3.2",
				Prompt:         "Hello world",
				PromptFile:     "",
				URL:            "http://localhost:11434/api/generate",
				Output:         "",
				DisableLoading: false,
				Silent:         false,
			},
			wantErr: false,
		},
		{
			name: "Prompt from file",
			args: []string{"cmd", "--prompt-file", promptFilePath},
			wantConfig: &Config{
				Model:          "llama3.2",
				Prompt:         "Hello from file",
				PromptFile:     promptFilePath,
				URL:            "http://localhost:11434/api/generate",
				Output:         "",
				DisableLoading: false,
				Silent:         false,
			},
			wantErr: false,
		},
		{
			name:           "Prompt file does not exist",
			args:           []string{"cmd", "--prompt-file", "nonexistent.txt"},
			wantErr:        true,
			wantErrMessage: "error reading prompt file 'nonexistent.txt': open nonexistent.txt: no such file or directory",
		},
		{
			name:           "Silent mode without output",
			args:           []string{"cmd", "--prompt=Hello", "--silent"},
			wantErr:        true,
			wantErrMessage: "the -silent flag requires the -output flag to be specified",
		},
		{
			name: "Silent mode with output",
			args: []string{"cmd", "--prompt=Hello", "--silent", "--output=result.txt"},
			wantConfig: &Config{
				Model:          "llama3.2",
				Prompt:         "Hello",
				PromptFile:     "",
				URL:            "http://localhost:11434/api/generate",
				Output:         "result.txt",
				DisableLoading: false,
				Silent:         true,
			},
			wantErr: false,
		},
		{
			name:     "Environment variables for model and URL",
			args:     []string{"cmd", "--prompt=Hello"},
			envModel: "env_model",
			envURL:   "http://env-url/api",
			wantConfig: &Config{
				Model:          "env_model",
				Prompt:         "Hello",
				PromptFile:     "",
				URL:            "http://env-url/api",
				Output:         "",
				DisableLoading: false,
				Silent:         false,
			},
			wantErr: false,
		},
		{
			name: "Disable loading and silent mode",
			args: []string{"cmd", "--prompt=Hello", "--no-loading", "--silent", "--output=result.txt"},
			wantConfig: &Config{
				Model:          "llama3.2",
				Prompt:         "Hello",
				PromptFile:     "",
				URL:            "http://localhost:11434/api/generate",
				Output:         "result.txt",
				DisableLoading: true,
				Silent:         true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save and restore os.Args and environment variables
			origArgs := os.Args
			origEnvModel := os.Getenv("NINO_MODEL")
			origEnvURL := os.Getenv("NINO_URL")

			defer func() {
				os.Args = origArgs
				if origEnvModel != "" {
					os.Setenv("NINO_MODEL", origEnvModel)
				} else {
					os.Unsetenv("NINO_MODEL")
				}
				if origEnvURL != "" {
					os.Setenv("NINO_URL", origEnvURL)
				} else {
					os.Unsetenv("NINO_URL")
				}
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError) // Reset flags
			}()

			// Set os.Args to the test case arguments
			os.Args = tt.args

			// Set environment variables if specified
			if tt.envModel != "" {
				os.Setenv("NINO_MODEL", tt.envModel)
			} else {
				os.Unsetenv("NINO_MODEL")
			}
			if tt.envURL != "" {
				os.Setenv("NINO_URL", tt.envURL)
			} else {
				os.Unsetenv("NINO_URL")
			}

			// Parse the arguments
			gotConfig, err := ParseArgs()
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseArgs() expected error but got none")
				} else if tt.wantErrMessage != "" && err.Error() != tt.wantErrMessage {
					t.Errorf("ParseArgs() error = '%v', want '%v'", err.Error(), tt.wantErrMessage)
				}
				return
			} else {
				if err != nil {
					t.Errorf("ParseArgs() unexpected error: %v", err)
					return
				}
			}

			// Check the obtained config
			if !reflect.DeepEqual(gotConfig, tt.wantConfig) {
				t.Errorf("ParseArgs() = %v, want %v", gotConfig, tt.wantConfig)
			}
		})
	}
}
