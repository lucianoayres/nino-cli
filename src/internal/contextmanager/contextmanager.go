package contextmanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// getDataDir returns the data directory, checking XDG_DATA_HOME first.
func getDataDir() (string, error) {
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("unable to determine home directory: %v", err)
		}
		dataDir = filepath.Join(homeDir, ".local", "share")
	}
	return dataDir, nil
}

// sanitizeModelName sanitizes the model name for use in file paths.
func sanitizeModelName(s string) string {
	// Replace any character that is not a letter, digit, or allowed punctuation with '_'
	// Allowed characters: letters, digits, '-', '_', and '.'
	re := regexp.MustCompile(`[^a-zA-Z0-9._-]+`)
	return re.ReplaceAllString(s, "_")
}

// SaveContext saves the context data for a given model.
func SaveContext(modelName string, context []int) error {
	dataDir, err := getDataDir()
	if err != nil {
		return err
	}

	modelName = sanitizeModelName(modelName)

	path := filepath.Join(dataDir, "nino", "models", modelName)

	// Create directories if they do not exist
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directories: %v", err)
	}

	contextFile := filepath.Join(path, "context.json")

	file, err := os.Create(contextFile)
	if err != nil {
		return fmt.Errorf("failed to create context file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(context)
	if err != nil {
		return fmt.Errorf("failed to encode context data: %v", err)
	}

	return nil
}

// LoadContext loads the context data for a given model.
func LoadContext(modelName string) ([]int, error) {
	dataDir, err := getDataDir()
	if err != nil {
		return nil, err
	}

	modelName = sanitizeModelName(modelName)

	path := filepath.Join(dataDir, "nino", "models", modelName, "context.json")

	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Context file does not exist
			return nil, nil
		}
		return nil, fmt.Errorf("failed to open context file: %v", err)
	}
	defer file.Close()

	var context []int
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&context)
	if err != nil {
		return nil, fmt.Errorf("failed to decode context data: %v", err)
	}

	return context, nil
}
