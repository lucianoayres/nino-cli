package processor

import (
	"encoding/json"
	"fmt"
	"io"
	"nino/internal/models"
)

// ProcessResponse reads and processes the response from the server
// It writes the response to the provided writer without altering the original formatting.
// It also handles saving context data when r.Done is true.
func ProcessResponse(body io.Reader, writer io.Writer, contextHandler func([]int) error) error {
	decoder := json.NewDecoder(body)

	for {
		var r models.ResponsePayload
		if err := decoder.Decode(&r); err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("failed to decode JSON response: %v", err)
		}

		if r.Response != "" {
			fmt.Fprint(writer, r.Response)
		}

		if r.Done {
			if len(r.Context) > 0 && contextHandler != nil {
				if err := contextHandler(r.Context); err != nil {
					return fmt.Errorf("failed to handle context: %v", err)
				}
			}
			break
		}
	}
	return nil
}
