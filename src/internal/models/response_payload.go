package models

// ResponsePayload represents the structure of the JSON response
type ResponsePayload struct {
    Model     string `json:"model"`
    CreatedAt string `json:"created_at"`
    Response  string `json:"response"`
    Done      bool   `json:"done"`
}
