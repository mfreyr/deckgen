package llm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mfreyr/deckgen/internal/models"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenAIClient_Generate(t *testing.T) {
	// The expected presentation that the mock server will return.
	expectedPresentation := models.Presentation{
		Job: models.Job{
			Role:    "Software Engineer",
			Company: "TestCorp",
		},
		ExperienceSummary: "A great summary.",
	}

	// Create a mock server to simulate the OpenAI API.
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// We expect a POST request to the chat completions endpoint.
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/chat/completions", r.URL.Path)

		// Create the mock response body.
		mockResponse := openai.ChatCompletionResponse{
			Choices: []openai.ChatCompletionChoice{
				{
					Message: openai.ChatCompletionMessage{
						Role:    openai.ChatMessageRoleAssistant,
						Content: mustMarshal(t, expectedPresentation),
					},
				},
			},
		}

		// Send the response.
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(mockResponse)
		require.NoError(t, err)
	}))
	defer mockServer.Close()

	// Configure the go-openai client to use the mock server.
	config := openai.DefaultConfig("fake-api-key")
	config.BaseURL = mockServer.URL + "/v1"
	client := openai.NewClientWithConfig(config)

	// Create our OpenAIClient with the mocked underlying client.
	llmClient := &OpenAIClient{client: client}

	// Call the Generate method.
	presentation, err := llmClient.Generate("some resume", "some job ad")

	// Assert the results.
	require.NoError(t, err)
	assert.Equal(t, expectedPresentation.Job.Role, presentation.Job.Role)
	assert.Equal(t, expectedPresentation.Job.Company, presentation.Job.Company)
	assert.Equal(t, expectedPresentation.ExperienceSummary, presentation.ExperienceSummary)
}

// mustMarshal is a helper to marshal a struct to JSON and fail the test on error.
func mustMarshal(t *testing.T, v interface{}) string {
	t.Helper()
	data, err := json.Marshal(v)
	require.NoError(t, err)
	return string(data)
}