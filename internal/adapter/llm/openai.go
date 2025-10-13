package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/mfreyr/deckgen/internal/config"
	"github.com/mfreyr/deckgen/internal/model"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/responses"
)

const (
	parseResumePromptTemplate = `
		**Objective:**
		Analyze the provided raw text from a resume file.
		Extract the information and structure it into a valid JSON object that adheres exactly to the provided JSON schema.

		**Instructions:**
		1. Parse the document to identify key sections like professional summary, work experience, skills, and certifications.
		2. Populate all fields of the JSON schema as accurately as possible.
		3. The output MUST be a single, valid JSON object. Do not include any text, markdown, or commentary outside of the JSON object.

		**Input Data (Raw Text from Resume):**
		---
		%s
	`
	parseJobAdPromptTemplate = `
		**Objective:**
		Analyze the provided raw text from a job advertisement.
		Extract the information and structure it into a valid JSON object that adheres exactly to the provided JSON schema.

		**Instructions:**
		1. Parse the document to identify key sections like job title, company name, responsibilities, and qualifications.
		2. Populate all fields of the JSON schema as accurately as possible.
		3. The output MUST be a single, valid JSON object. Do not include any text, markdown, or commentary outside of the JSON object.

		**Input Data (Raw Text from Job Ad):**
		---
		%s
	`
	adaptResumePromptTemplate = `
		**Objective:**
		Analyze the provided Job Advertisement and one or more candidate resumes.
		Generate a new, adapted resume in JSON format that highlights the candidate's most relevant skills and experiences for this specific job.

		**Instructions:**
		1.  Carefully read the Job Advertisement to understand the key requirements, skills, and responsibilities.
		2.  Thoroughly review all provided candidate resumes to understand the candidate's background, skills, and accomplishments.
		3.  Synthesize this information to create compelling, concise, and action-oriented content for a new, adapted resume.
		4.  The output MUST be a single, valid JSON object that adheres exactly to the schema provided for the adapted resume.

		**Input Data:**

		--- Job Advertisement ---
		%s

		--- Candidate Resumes ---
		%s
	`
)

// OpenAIProvider implements the service.LLMProvider interface for OpenAI's models.
type OpenAIProvider struct {
	client    *openai.Client
	modelName string
}

// NewOpenAIProvider initializes and returns a new OpenAIProvider using the given config.
func NewOpenAIProvider(cfg config.LLMProviderConfig) (*OpenAIProvider, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("openAI API key is required")
	}
	if cfg.Model == "" {
		return nil, fmt.Errorf("openAI model name is required in config")
	}

	client := openai.NewClient(option.WithAPIKey(cfg.APIKey))

	return &OpenAIProvider{
		client:    &client,
		modelName: cfg.Model,
	}, nil
}

// ParseResume uses an LLM to parse a file into a structured CandidateResume.
func (p *OpenAIProvider) ParseResume(ctx context.Context, file model.File) (model.CandidateResume, error) {
	var resume model.CandidateResume
	prompt := fmt.Sprintf(parseResumePromptTemplate, string(file.Content))

	fileParam := openai.FileNewParams{
		File:    openai.File(bytes.NewReader(file.Content), "resume.pdf", "application/pdf"),
		Purpose: openai.FilePurposeUserData,
	}
	storedFile, err := p.client.Files.New(ctx, fileParam)
	if err != nil {
		return model.CandidateResume{}, fmt.Errorf("error uploading file to OpenAI: %w", err)
	}

	params := responses.ResponseNewParams{
		Model: openai.ChatModel(p.modelName),
		Text: responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfJSONSchema: &responses.ResponseFormatTextJSONSchemaConfigParam{
					Name:        "Parsed Resume",
					Description: openai.String("Structured json resume parsed from a file"),
					Schema:      GenerateSchema[model.CandidateResume](),
					Strict:      openai.Bool(true),
				},
			},
		},
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: responses.ResponseInputParam{
				responses.ResponseInputItemParamOfMessage(
					responses.ResponseInputMessageContentListParam{
						responses.ResponseInputContentUnionParam{
							OfInputFile: &responses.ResponseInputFileParam{
								FileID: openai.String(storedFile.ID),
							},
						},
						responses.ResponseInputContentUnionParam{
							OfInputText: &responses.ResponseInputTextParam{
								Text: prompt,
							},
						},
					},
					"user",
				),
			},
		},
	}

	rawJSON, err := p.executeRequest(ctx, params)
	if err != nil {
		return resume, err
	}

	if err := json.Unmarshal([]byte(rawJSON), &resume); err != nil {
		log.Printf("Failed to unmarshal JSON from OpenAI for ParseResume. Raw response:\n%s", rawJSON)
		return resume, fmt.Errorf("failed to unmarshal JSON from OpenAI: %w", err)
	}

	return resume, nil
}

// ParseJobAd uses an LLM to parse a file into a structured JobAd.
func (p *OpenAIProvider) ParseJobAd(ctx context.Context, file model.File) (model.JobAd, error) {
	var jobAd model.JobAd
	prompt := fmt.Sprintf(parseJobAdPromptTemplate, string(file.Content))

	fileParam := openai.FileNewParams{
		File:    openai.File(bytes.NewReader(file.Content), "job_ad.pdf", "application/pdf"),
		Purpose: openai.FilePurposeUserData,
	}
	storedFile, err := p.client.Files.New(ctx, fileParam)
	if err != nil {
		return model.JobAd{}, fmt.Errorf("error uploading file to OpenAI: %w", err)
	}

	params := responses.ResponseNewParams{
		Model: openai.ChatModel(p.modelName),
		Text: responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfJSONSchema: &responses.ResponseFormatTextJSONSchemaConfigParam{
					Name:        "Parsed Job Ad",
					Description: openai.String("Structured json of job ad parsed from a file"),
					Schema:      GenerateSchema[model.CandidateResume](),
					Strict:      openai.Bool(true),
				},
			},
		},
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: responses.ResponseInputParam{
				responses.ResponseInputItemParamOfMessage(
					responses.ResponseInputMessageContentListParam{
						responses.ResponseInputContentUnionParam{
							OfInputFile: &responses.ResponseInputFileParam{
								FileID: openai.String(storedFile.ID),
							},
						},
						responses.ResponseInputContentUnionParam{
							OfInputText: &responses.ResponseInputTextParam{
								Text: prompt,
							},
						},
					},
					"user",
				),
			},
		},
	}
	rawJSON, err := p.executeRequest(ctx, params)
	if err != nil {
		return jobAd, err
	}

	if err := json.Unmarshal([]byte(rawJSON), &jobAd); err != nil {
		log.Printf("Failed to unmarshal JSON from OpenAI for ParseJobAd. Raw response:\n%s", rawJSON)
		return jobAd, fmt.Errorf("failed to unmarshal JSON from OpenAI: %w", err)
	}

	return jobAd, nil
}

// AdaptResume uses an LLM to tailor existing resumes for a specific job ad.
func (p *OpenAIProvider) AdaptResume(ctx context.Context, jobAd model.JobAd, resumes []model.CandidateResume) (model.CandidateAdaptedResume, error) {
	var adaptedResume model.CandidateAdaptedResume

	var resumeBuilder strings.Builder
	for i, resume := range resumes {
		resumeBytes, err := json.Marshal(resume)
		if err != nil {
			return adaptedResume, fmt.Errorf("failed to marshal resume ID %d to JSON: %w", resume.ID, err)
		}
		resumeBuilder.WriteString(fmt.Sprintf("\n--- Candidate Resume %d ---\n%s", i+1, string(resumeBytes)))
	}

	jobAdBytes, err := json.Marshal(jobAd)
	if err != nil {
		return adaptedResume, fmt.Errorf("failed to marshal job ad to JSON: %w", err)
	}

	prompt := fmt.Sprintf(adaptResumePromptTemplate, string(jobAdBytes), resumeBuilder.String())

	params := responses.ResponseNewParams{
		Model: openai.ChatModel(p.modelName),
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(prompt),
		},
	}
	rawJSON, err := p.executeRequest(ctx, params)
	if err != nil {
		return adaptedResume, err
	}

	if err := json.Unmarshal([]byte(rawJSON), &adaptedResume); err != nil {
		log.Printf("Failed to unmarshal JSON from OpenAI for AdaptResume. Raw response:\n%s", rawJSON)
		return adaptedResume, fmt.Errorf("failed to unmarshal JSON from OpenAI: %w", err)
	}

	return adaptedResume, nil
}

// executeRequest is a helper function to run the chat completion and handle the response.
func (p *OpenAIProvider) executeRequest(ctx context.Context, params responses.ResponseNewParams) (string, error) {
	resp, err := p.client.Responses.New(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to create responses with OpenAI: %w", err)
	}
	return resp.OutputText(), nil
}
