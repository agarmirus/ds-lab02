package models

type ValidationErrorResponse struct {
	Message string             `json:"message"`
	Errors  []ErrorDiscription `json:"errors"`
}
