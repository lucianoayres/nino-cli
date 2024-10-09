// models/models.go
package models

// ResponsePayload represents the structure of each JSON object in the response stream.
type ResponsePayload struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

// RequestPayload represents the payload sent in the HTTP request.
type RequestPayload struct {
	Model  string   `json:"model"`
	Prompt string   `json:"prompt"`
	Images []string `json:"images"` // New field for images in base64
}
