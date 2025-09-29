package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"github.com/mfreyr/deckgen/internal/models"

	"github.com/sashabaranov/go-openai"
)

// OpenAIClient implements the LLM interface using the OpenAI API.
type OpenAIClient struct {
	client *openai.Client
}

// NewOpenAIClient creates a new client for interacting with the OpenAI API.
// It expects the OPENAI_API_KEY environment variable to be set.
func NewOpenAIClient() *OpenAIClient {
	// The NewClient function without an API key will default to reading
	// the key from the OPENAI_API_KEY environment variable.
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	return &OpenAIClient{client: client}
}

// Generate uses the OpenAI API to adapt a resume for a specific job advertisement.
func (c *OpenAIClient) Generate(resumeText string, jobAdText string) (*models.Presentation, error) {
	prompt := buildPrompt(resumeText, jobAdText)

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are an expert career assistant. Your task is to analyze a resume and a job description, and then generate a professional presentation in JSON format. The JSON output must strictly follow the provided structure. Do not include any introductory text or markdown formatting around the JSON.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("error calling OpenAI API: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response choices from OpenAI")
	}

	var presentation models.Presentation
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &presentation)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling OpenAI response: %w", err)
	}

	return &presentation, nil
}

// buildPrompt constructs the detailed prompt for the LLM.
func buildPrompt(resumeText, jobAdText string) string {
	return fmt.Sprintf(`
Here is the resume content:
---
%s
---

Here is the job advertisement:
---
%s
---

Based on the resume and job ad, please generate a JSON object with the following structure:

{
  "Job": {
    "Role": "The role from the job ad",
    "Company": "The company from the job ad"
  },
  "ExperienceSummary": "A 2-3 sentence summary of the candidate's most relevant experience for this specific job.",
  "Experiences": [
    {
      "Firm": "Name of the company",
      "Role": "Job title",
      "Date": "Dates of employment",
      "Description": "A brief description of the role, tailored to highlight relevance to the job ad.",
      "Tasks": ["Task 1 relevant to the job ad", "Task 2 relevant to the job ad"],
      "Technologies": ["Tech 1", "Tech 2"]
    }
  ],
  "CommercialProposition": {
    "Role": "Proposed role for the candidate",
    "ConsultantName": "Candidate's Name (extract from resume)",
    "Localisation": "Job location (from job ad)",
    "Availability": "Candidate's availability (if mentioned, otherwise state 'To be discussed')"
  },
  "ContactCommercial": {
    "Name": "Your Name (as the assistant)",
    "Email": "your.email@example.com",
    "Phone": "your-phone-number"
  }
}

Please fill in the details based *only* on the provided resume and job advertisement.
`, resumeText, jobAdText)
}