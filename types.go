// Package jsonrepair provides functionality to repair invalid JSON documents.
// This is a Go port of the original jsonrepair project by Jos de Jong.
// Original project: https://github.com/josdejong/jsonrepair
// Licensed under the ISC License
package jsonrepair

import (
	"fmt"
	"strings"
)

// JSONRepairError represents an error that occurred during JSON repair
type JSONRepairError struct {
	Message  string
	Position int
}

// Error implements the error interface
func (e *JSONRepairError) Error() string {
	return fmt.Sprintf("%s at position %d", e.Message, e.Position)
}

// NewJSONRepairError creates a new JSONRepairError
func NewJSONRepairError(message string, position int) *JSONRepairError {
	return &JSONRepairError{
		Message:  message,
		Position: position,
	}
}

// Parser represents a JSON repair parser
type Parser struct {
	text   string          // Input text to parse
	output strings.Builder // Output buffer for repaired JSON
	i      int             // Current position index in text
}

// NewParser creates a new Parser instance
func NewParser(text string) *Parser {
	return &Parser{
		text: text,
		i:    0,
	}
}
