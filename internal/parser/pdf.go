package parser

import (
	"bytes"
	"io"
	"os/exec"
)

// PDFParser implements the Parser interface for PDF files.
// It uses the external `pdftotext` command-line utility to extract text.
// This approach is robust and avoids the need for a pure Go PDF library.
type PDFParser struct{}

// NewPDFParser creates a new PDFParser.
func NewPDFParser() *PDFParser {
	return &PDFParser{}
}

// Parse extracts text from a PDF file by shelling out to the `pdftotext` utility.
// It reads the PDF content from the provided reader, pipes it to the command's stdin,
// and captures the extracted text from stdout.
// The arguments `"-" "-"` instruct `pdftotext` to read the PDF from stdin and write the text to stdout.
func (p *PDFParser) Parse(reader io.Reader) (string, error) {
	// Create the command to execute `pdftotext`.
	// The arguments are:
	// -: read from stdin
	// -: write to stdout
	cmd := exec.Command("pdftotext", "-", "-")

	// Set the command's standard input to the provided reader.
	cmd.Stdin = reader

	// Create a buffer to capture the command's standard output.
	var out bytes.Buffer
	cmd.Stdout = &out

	// Run the command.
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Return the captured output as a string.
	return out.String(), nil
}