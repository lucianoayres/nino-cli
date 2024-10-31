package models

// ResponsePayload represents the structure of each JSON object in the response stream.
type ResponsePayload struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
	Context   []int  `json:"context"`
}

// RequestPayload represents the payload sent in the HTTP request.
type RequestPayload struct {
	Model      string   `json:"model"`
	Prompt     string   `json:"prompt"`
	Images     []string `json:"images"` // New field for images in base64
	Format     string   `json:"format"`
	Stream     bool     `json:"stream"`
	Keep_Alive string   `json:"keep_alive,omitempty"`
	Context    []int    `json:"context,omitempty"`
}
