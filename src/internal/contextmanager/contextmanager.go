package contextmanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/lucianoayres/nino-cli/internal/logger"
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
	log := logger.GetLogger(true) // Assuming logger is already initialized in main
	log.Info("Saving context data for model: %s", modelName)

	dataDir, err := getDataDir()
	if err != nil {
		log.Error("Failed to get data directory: %v", err)
		return err
	}

	modelName = sanitizeModelName(modelName)

	path := filepath.Join(dataDir, "nino", "models", modelName)

	// Create directories if they do not exist
	log.Info("Creating directories if not exist: %s", path)
	err = os.MkdirAll(path, 0755)
	if err != nil {
		log.Error("Failed to create directories: %v", err)
		return fmt.Errorf("failed to create directories: %v", err)
	}

	contextFile := filepath.Join(path, "context.json")

	log.Info("Creating context file: %s", contextFile)
	file, err := os.Create(contextFile)
	if err != nil {
		log.Error("Failed to create context file: %v", err)
		return fmt.Errorf("failed to create context file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	log.Info("Encoding context data to JSON")
	err = encoder.Encode(context)
	if err != nil {
		log.Error("Failed to encode context data: %v", err)
		return fmt.Errorf("failed to encode context data: %v", err)
	}

	log.Info("Context data saved successfully")
	return nil
}

// LoadContext loads the context data for a given model.
func LoadContext(modelName string) ([]int, error) {
	log := logger.GetLogger(true) // Assuming logger is already initialized in main
	log.Info("Loading context data for model: %s", modelName)

	dataDir, err := getDataDir()
	if err != nil {
		log.Error("Failed to get data directory: %v", err)
		return nil, err
	}

	modelName = sanitizeModelName(modelName)

	path := filepath.Join(dataDir, "nino", "models", modelName, "context.json")

	log.Info("Opening context file: %s", path)
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Context file does not exist
			log.Info("Context file does not exist: %s", path)
			return nil, nil
		}
		log.Error("Failed to open context file: %v", err)
		return nil, fmt.Errorf("failed to open context file: %v", err)
	}
	defer file.Close()

	var context []int
	decoder := json.NewDecoder(file)
	log.Info("Decoding context data from JSON")
	err = decoder.Decode(&context)
	if err != nil {
		log.Error("Failed to decode context data: %v", err)
		return nil, fmt.Errorf("failed to decode context data: %v", err)
	}

	log.Info("Context data loaded successfully")
	return context, nil
}
