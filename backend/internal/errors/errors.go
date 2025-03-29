package errors

import "errors"

// Custom errors
var (
	ErrInvalidInput         = errors.New("invalid input")
	ErrProjectNotFound      = errors.New("project not found")
	ErrProjectAlreadyExists = errors.New("project already exists")
	ErrIssueNotFound        = errors.New("issue not found")
)
