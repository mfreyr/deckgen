package parser

import "io"

// Parser defines the interface for document parsers.
// It allows for different implementations to handle various document types (PDF, DOCX, etc.).
type Parser interface {
	// Parse reads content from an io.Reader, extracts the text, and returns it as a string.
	Parse(reader io.Reader) (string, error)
}