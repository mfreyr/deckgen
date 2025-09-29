package api

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/mfreyr/deckgen/internal/llm"
	"github.com/mfreyr/deckgen/internal/parser"
)

// GenerateHandler is the HTTP handler for the /generate endpoint.
// It orchestrates the process of parsing a resume, calling the LLM,
// and returning the structured presentation.
type GenerateHandler struct {
	parser parser.Parser
	llm    llm.LLM
}

// NewGenerateHandler creates a new GenerateHandler with its dependencies.
func NewGenerateHandler(p parser.Parser, l llm.LLM) *GenerateHandler {
	return &GenerateHandler{
		parser: p,
		llm:    l,
	}
}

// ServeHTTP handles the incoming HTTP request.
func (h *GenerateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form, with a 10MB max memory limit.
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
		return
	}

	// Get the job ad from the form values.
	jobAd := r.FormValue("jobAd")
	if jobAd == "" {
		http.Error(w, "Missing 'jobAd' form field", http.StatusBadRequest)
		return
	}

	// Get the resume file from the form.
	file, _, err := r.FormFile("resume")
	if err != nil {
		http.Error(w, "Error retrieving 'resume' file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Use the injected parser to extract text from the file.
	resumeText, err := h.parser.Parse(file)
	if err != nil {
		log.Printf("Error parsing file: %v", err)
		http.Error(w, "Error parsing resume file", http.StatusInternalServerError)
		return
	}

	// Use the injected LLM client to generate the presentation.
	presentation, err := h.llm.Generate(resumeText, jobAd)
	if err != nil {
		log.Printf("Error generating presentation: %v", err)
		http.Error(w, "Error generating presentation from LLM", http.StatusInternalServerError)
		return
	}

	// Send the successful JSON response.
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(presentation); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}