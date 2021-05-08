package controllers

type NoError struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
}
