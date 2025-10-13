package source

import (
	"github.com/google/uuid"
)

// TextSource is a concrete implementation of the KnowledgeSource interface
// for raw text provided by the user (e.g., from a textarea).
type TextSource struct {
	id      string
	name    string
	content string
}

// NewTextSource creates a new TextSource.
func NewTextSource(sourceName, textContent string) *TextSource {
	return &TextSource{
		id:      uuid.New().String(),
		name:    sourceName,
		content: textContent,
	}
}

// GetID returns the unique identifier for this text source.
func (s *TextSource) GetID() string {
	return s.id
}

// GetName returns the descriptive name for this text source.
func (s *TextSource) GetName() string {
	return s.name
}

// IsFile for a TextSource always returns false.
func (s *TextSource) IsFile() bool {
	return false
}

// GetParsedContent for a TextSource simply returns the raw text content.
// No parsing is required.
func (s *TextSource) GetParsedContent() (string, error) {
	return s.content, nil
}
