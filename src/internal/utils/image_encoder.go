package utils

import (
	"encoding/base64"
	"fmt"
	"os"
)

// ReadImagesAsBase64 reads image files from the given paths and returns their base64-encoded strings.
func ReadImagesAsBase64(paths []string) ([]string, error) {
	var images []string
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("error reading image file '%s': %v", path, err)
		}
		encoded := base64.StdEncoding.EncodeToString(data)
		images = append(images, encoded)
	}
	return images, nil
}