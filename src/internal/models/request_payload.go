package models

// RequestPayload represents the JSON payload sent in the request
type RequestPayload struct {
    Model  string `json:"model"`
    Prompt string `json:"prompt"`
}
