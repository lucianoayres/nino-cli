package utils

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"github.com/lucianoayres/nino-cli/internal/logger"
)

func TestReadImagesAsBase64(t *testing.T) {
	// Initialize logger for tests
	logger.GetLogger(true)

	t.Run("Valid image files", func(t *testing.T) {
		// Create a temporary directory
		tempDir := t.TempDir()

		// Create a temporary image file
		imagePath := filepath.Join(tempDir, "test_image.jpg")
		imageContent := []byte{0xFF, 0xD8, 0xFF, 0xE0} // JPEG file header bytes
		err := os.WriteFile(imagePath, imageContent, 0644)
		if err != nil {
			t.Fatalf("Failed to create temporary image file: %v", err)
		}

		// Expected base64 encoding of the imageContent
		expectedBase64 := base64.StdEncoding.EncodeToString(imageContent)

		// Call the function under test
		images, err := ReadImagesAsBase64([]string{imagePath})
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if len(images) != 1 {
			t.Fatalf("Expected 1 image, got: %d", len(images))
		}

		if images[0] != expectedBase64 {
			t.Errorf("Base64 output does not match expected value.\nExpected: %s\nGot: %s", expectedBase64, images[0])
		}
	})

	t.Run("Invalid file path", func(t *testing.T) {
		// Call the function with an invalid path
		_, err := ReadImagesAsBase64([]string{"/invalid/path/image.jpg"})
		if err == nil {
			t.Fatalf("Expected an error for invalid file path, got nil")
		}
	})

	t.Run("Empty input", func(t *testing.T) {
		// Call the function with an empty slice
		images, err := ReadImagesAsBase64([]string{})
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if len(images) != 0 {
			t.Errorf("Expected 0 images, got: %d", len(images))
		}
	})
}
