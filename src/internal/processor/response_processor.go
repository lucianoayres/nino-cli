package processor

import (
	"encoding/json"
	"fmt"
	"io"

	"gollama/internal/models"
)

// ProcessResponse reads and processes the response from the server
func ProcessResponse(body io.Reader) error {
    decoder := json.NewDecoder(body)

    for {
        var r models.ResponsePayload
        if err := decoder.Decode(&r); err == io.EOF {
            break
        } else if err != nil {
            return fmt.Errorf("failed to decode JSON response: %v", err)
        }

        fmt.Print(r.Response)

        if r.Done {
            break
        }
    }
    return nil
}
