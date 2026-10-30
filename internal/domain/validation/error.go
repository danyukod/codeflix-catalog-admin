package validation

import "strings"

// ValidationError contains all validation failures raised by a domain operation.
type ValidationError struct {
	Messages []string
}

func NewValidationError(errors []Error) *ValidationError {
	messages := make([]string, 0, len(errors))
	for _, err := range errors {
		messages = append(messages, err.Message)
	}
	return &ValidationError{Messages: messages}
}

func (e *ValidationError) Error() string {
	return "validation failed: " + strings.Join(e.Messages, "; ")
}
