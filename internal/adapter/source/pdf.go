package source

import (
	"bytes"
	"fmt"

	"code.sajari.com/docconv/v2"
	"github.com/google/uuid"
)

// PDFSource is a concrete implementation of the KnowledgeSource interface for PDF files.
type PDFSource struct {
	id      string
	name    string
	content []byte
}

// NewPDFSource creates a new PDFSource from a filename and its byte content.
func NewPDFSource(fileName string, fileContent []byte) *PDFSource {
	return &PDFSource{
		id:      uuid.New().String(),
		name:    fileName,
		content: fileContent,
	}
}

// GetID returns the unique identifier for this PDF source.
func (s *PDFSource) GetID() string {
	return s.id
}

// GetName returns the original filename of the PDF.
func (s *PDFSource) GetName() string {
	return s.name
}

// IsFile for a PDFSource always returns true.
func (s *PDFSource) IsFile() bool {
	return true
}

// GetParsedContent is responsible for extracting text from the PDF content.
// In a real implementation, this method would use a PDF parsing library.
func (s *PDFSource) GetParsedContent() (string, error) {
	r := bytes.NewReader(s.content)
	res, _, err := docconv.ConvertPDF(r)
	if err != nil {
		return "", fmt.Errorf("counld not read pdf: %w", err)
	}
	return res, nil
}
