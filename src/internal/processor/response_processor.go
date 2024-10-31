package processor

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/lucianoayres/nino-cli/internal/logger"
	"github.com/lucianoayres/nino-cli/internal/models"
)

// ProcessResponse reads and processes the response from the server
// It writes the response to the provided writer without altering the original formatting.
// It also handles saving context data when r.Done is true.
func ProcessResponse(body io.Reader, writer io.Writer, contextHandler func([]int) error) error {
	log := logger.GetLogger(true) // Assuming logger is already initialized in main
	log.Info("Starting to process response")
	decoder := json.NewDecoder(body)

	for {
		var r models.ResponsePayload
		if err := decoder.Decode(&r); err == io.EOF {
			log.Info("End of response stream")
			break
		} else if err != nil {
			log.Error("JSON decoding error: %v", err)
			return fmt.Errorf("failed to decode JSON response: %v", err)
		}

		log.Info("Received ResponsePayload: Model=%s, CreatedAt=%s, Done=%v", r.Model, r.CreatedAt, r.Done)

		if r.Response != "" {
			log.Info("Writing response to writer: %s", r.Response)
			fmt.Fprint(writer, r.Response)
		}

		if r.Done {
			log.Info("Response marked as done")
			if len(r.Context) > 0 && contextHandler != nil {
				log.Info("Handling context data: %v", r.Context)
				if err := contextHandler(r.Context); err != nil {
					log.Error("Context handler error: %v", err)
					return fmt.Errorf("failed to handle context: %v", err)
				}
			}
			break
		}
	}
	log.Info("Finished processing response")
	return nil
}
