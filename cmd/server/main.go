package main

import (
	"log"
	"net/http"
	"github.com/mfreyr/deckgen/internal/api"
	"github.com/mfreyr/deckgen/internal/llm"
	"github.com/mfreyr/deckgen/internal/parser"
)

func main() {
	// 1. Initialize dependencies.
	// We create a new PDF parser and a new OpenAI client.
	// These are the concrete implementations of our interfaces.
	pdfParser := parser.NewPDFParser()
	openaiClient := llm.NewOpenAIClient()

	// 2. Inject dependencies into the handler.
	// The handler now has access to the parser and LLM client
	// without being tightly coupled to the specific implementations.
	generateHandler := api.NewGenerateHandler(pdfParser, openaiClient)

	// 3. Set up the HTTP router.
	// We create a new ServeMux, which is Go's standard HTTP request router.
	mux := http.NewServeMux()
	mux.Handle("/generate", generateHandler)

	// 4. Start the server.
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}