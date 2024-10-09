// processor/processor.go
package processor

import (
	"encoding/json"
	"fmt"
	"io"
	"nino/internal/models"
)

// ProcessResponse reads and processes the response from the server
// It writes the response to the provided writer without altering the original formatting.
func ProcessResponse(body io.Reader, writer io.Writer) error {
	decoder := json.NewDecoder(body)

	for {
		var r models.ResponsePayload
		if err := decoder.Decode(&r); err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("failed to decode JSON response: %v", err)
		}

		if r.Response == "" {
			continue // Skip empty responses
		}

		fmt.Fprint(writer, r.Response)

		if r.Done {
			break
		}
	}
	return nil
}
