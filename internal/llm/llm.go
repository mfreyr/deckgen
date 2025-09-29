package llm

import "github.com/mfreyr/deckgen/internal/models"

// LLM defines the interface for language model clients.
// This abstraction allows for swapping different LLM providers (e.g., OpenAI, Anthropic)
// without changing the core application logic.
type LLM interface {
	// Generate takes the text from a resume and a job ad, sends it to the language model,
	// and returns a structured Presentation object.
	Generate(resumeText string, jobAdText string) (*models.Presentation, error)
}