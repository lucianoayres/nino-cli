package utils

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/lucianoayres/nino-cli/internal/logger"
)

// ReadImagesAsBase64 reads image files from the given paths and returns their base64-encoded strings.
func ReadImagesAsBase64(paths []string) ([]string, error) {
	log := logger.GetLogger(true) // Assuming logger is already initialized in main
	log.Info("Reading %d image(s) for encoding", len(paths))

	var images []string
	for _, path := range paths {
		log.Info("Reading image file: %s", path)
		data, err := os.ReadFile(path)
		if err != nil {
			log.Error("Error reading image file '%s': %v", path, err)
			return nil, fmt.Errorf("error reading image file '%s': %v", path, err)
		}
		encoded := base64.StdEncoding.EncodeToString(data)
		log.Info("Image file '%s' encoded successfully", path)
		images = append(images, encoded)
	}
	log.Info("All images encoded successfully")
	return images, nil
}
