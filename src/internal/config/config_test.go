package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestParseArgs(t *testing.T) {
	// Create a temporary prompt file for testing
	tmpDir := t.TempDir()
	promptFilePath := filepath.Join(tmpDir, "prompt.txt")
	err := os.WriteFile(promptFilePath, []byte("Hello from file"), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary prompt file: %v", err)
	}

	// Create temporary image files for testing
	imageFilePath1 := filepath.Join(tmpDir, "image1.jpg")
	err = os.WriteFile(imageFilePath1, []byte{0xFF, 0xD8, 0xFF, 0xE0}, 0644) // JPEG header bytes
	if err != nil {
		t.Fatalf("Failed to create temporary image file 1: %v", err)
	}

	imageFilePath2 := filepath.Join(tmpDir, "image2.jpg")
	err = os.WriteFile(imageFilePath2, []byte{0xFF, 0xD8, 0xFF, 0xE1}, 0644) // JPEG header bytes
	if err != nil {
		t.Fatalf("Failed to create temporary image file 2: %v", err)
	}

	tests := []struct {
		name            string
		args            []string
		envModel        string
		envURL          string
		envSystemPrompt string
		wantConfig      *Config
		wantErr         bool
		wantErrMessage  string
	}{
		{
			name: "Valid arguments with long flags",
			args: []string{"cmd", "--model=llama3.1", "--prompt=Hello", "--url=http://localhost:11434/api/generate", "--output=result.txt", "--format=json", "--no-stream"},
			wantConfig: &Config{
				Model:          "llama3.1",
				Prompt:         "Hello",
				PromptFile:     "",
				URL:            "http://localhost:11434/api/generate",
				Output:         "result.txt",
				DisableLoading: false,
				Silent:         false,
				ImagePaths:     []string{},
				Format:         "json",
				Stream:         false,
			},
			wantErr: false,
		},
		{
			name:            "Environment variables including system prompt",
			args:            []string{"cmd", "--prompt=Hello", "--format=json"},
			envModel:        "env_model",
			envURL:          "http://env-url/api",
			envSystemPrompt: "System prompt:",
			wantConfig: &Config{
				Model:          "env_model",
				Prompt:         strings.TrimSpace("System prompt: Hello"),
				PromptFile:     "",
				URL:            "http://env-url/api",
				Output:         "",
				DisableLoading: false,
				Silent:         false,
				ImagePaths:     []string{},
				Format:         "json",
				Stream:         true,
			},
			wantErr: false,
		},
		{
			name: "Valid arguments with image files",
			args: []string{"cmd", "--prompt=Hello", "--image", imageFilePath1, "--image", imageFilePath2, "--format=json", "--no-stream"},
			wantConfig: &Config{
				Model:          "llama3.2",
				Prompt:         strings.TrimSpace("Hello " + imageFilePath1 + " " + imageFilePath2),
				PromptFile:     "",
				URL:            "http://localhost:11434/api/generate",
				Output:         "",
				DisableLoading: false,
				Silent:         false,
				ImagePaths:     []string{imageFilePath1, imageFilePath2},
				Format:         "json",
				Stream:         false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save and restore os.Args, environment variables, and flag.CommandLine
			origArgs := os.Args
			origEnvModel := os.Getenv("NINO_MODEL")
			origEnvURL := os.Getenv("NINO_URL")
			origEnvSystemPrompt := os.Getenv("NINO_SYSTEM_PROMPT")
			origFlagCommandLine := flag.CommandLine

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
				if origEnvSystemPrompt != "" {
					os.Setenv("NINO_SYSTEM_PROMPT", origEnvSystemPrompt)
				} else {
					os.Unsetenv("NINO_SYSTEM_PROMPT")
				}
				flag.CommandLine = origFlagCommandLine
			}()

			// Reset flags before each test
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			flag.CommandLine.Usage = func() {
				fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
				flag.CommandLine.PrintDefaults()
			}

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
			if tt.envSystemPrompt != "" {
				os.Setenv("NINO_SYSTEM_PROMPT", tt.envSystemPrompt)
			} else {
				os.Unsetenv("NINO_SYSTEM_PROMPT")
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