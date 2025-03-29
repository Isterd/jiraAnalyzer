package models

import "fmt"

// NotFoundError represents an error when a resource is not found
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", e.Message)
}

// InvalidInputError represents an error caused by invalid user input
type InvalidInputError struct {
	Message string
}

func (e *InvalidInputError) Error() string {
	return fmt.Sprintf("invalid input: %s", e.Message)
}
