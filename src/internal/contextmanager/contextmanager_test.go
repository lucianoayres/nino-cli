package contextmanager

import (
	"os"
	"reflect"
	"testing"

	"github.com/lucianoayres/nino-cli/internal/logger"
)

// TestSaveAndLoadContext tests the SaveContext and LoadContext functions.
func TestSaveAndLoadContext(t *testing.T) {
	// Initialize logger for tests
	logger.GetLogger(true)

	// Create a temporary directory to act as XDG_DATA_HOME
	tmpDir, err := os.MkdirTemp("", "nino_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Set XDG_DATA_HOME to the temporary directory
	os.Setenv("XDG_DATA_HOME", tmpDir)
	defer os.Unsetenv("XDG_DATA_HOME")

	modelName := "test-model"
	contextData := []int{1, 2, 3, 4, 5}

	// Save the context
	err = SaveContext(modelName, contextData)
	if err != nil {
		t.Errorf("SaveContext returned error: %v", err)
	}

	// Load the context
	loadedContext, err := LoadContext(modelName)
	if err != nil {
		t.Errorf("LoadContext returned error: %v", err)
	}

	// Check if the loaded context matches the saved context
	if !reflect.DeepEqual(loadedContext, contextData) {
		t.Errorf("Loaded context does not match saved context.\nExpected: %v\nGot: %v", contextData, loadedContext)
	}

	// Verify that the context file is overwritten on subsequent saves
	newContextData := []int{6, 7, 8}
	err = SaveContext(modelName, newContextData)
	if err != nil {
		t.Errorf("SaveContext returned error on second save: %v", err)
	}

	loadedContext, err = LoadContext(modelName)
	if err != nil {
		t.Errorf("LoadContext returned error after second save: %v", err)
	}

	if !reflect.DeepEqual(loadedContext, newContextData) {
		t.Errorf("Loaded context does not match updated context.\nExpected: %v\nGot: %v", newContextData, loadedContext)
	}
}

// TestLoadContextNoFile tests that LoadContext returns nil when the context file does not exist.
func TestLoadContextNoFile(t *testing.T) {
	// Initialize logger for tests
	logger.GetLogger(true)

	// Create a temporary directory to act as XDG_DATA_HOME
	tmpDir, err := os.MkdirTemp("", "nino_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Set XDG_DATA_HOME to the temporary directory
	os.Setenv("XDG_DATA_HOME", tmpDir)
	defer os.Unsetenv("XDG_DATA_HOME")

	modelName := "nonexistent-model"

	// Attempt to load context for a model that has no saved context
	loadedContext, err := LoadContext(modelName)
	if err != nil {
		t.Errorf("LoadContext returned error: %v", err)
	}

	if loadedContext != nil {
		t.Errorf("Expected nil context for nonexistent model, got: %v", loadedContext)
	}
}

// TestSanitizeModelName tests the sanitizeModelName function with various inputs.
func TestSanitizeModelName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"model-name", "model-name"},
		{"model name", "model_name"},
		{"model@name", "model_name"},
		{"model$name", "model_name"},
		{"model/name", "model_name"},
		{"model\\name", "model_name"},
		{"model:name", "model_name"},
		{"model.name", "model.name"},
		{"model-name_version1", "model-name_version1"},
		{"model*name?", "model_name_"},
		{"ModelName123", "ModelName123"},
		{"", ""},
		{"   ", "_"},
		{"model\nname", "model_name"},
	}

	for _, tt := range tests {
		result := sanitizeModelName(tt.input)
		if result != tt.expected {
			t.Errorf("sanitizeModelName(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}
