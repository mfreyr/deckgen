package service

import (
	"context"
	"fmt"

	"github.com/mfreyr/deckgen/internal/model"
)

type ResumeRepository interface {
	SaveAdaptedResume(ctx context.Context, adaptedResume model.CandidateAdaptedResume) (model.CandidateAdaptedResume, error)
	GetAdaptedResume(ctx context.Context, adaptedResumeID int) (model.CandidateAdaptedResume, error)
	ListAdaptedResumes(ctx context.Context) ([]model.CandidateAdaptedResume, error)
	UpdateAdaptedResume(ctx context.Context, adaptedResume model.CandidateAdaptedResume) (model.CandidateAdaptedResume, error)
	DeleteAdaptedResume(ctx context.Context, adaptedResumeID int) error

	SaveResume(ctx context.Context, resume model.CandidateResume) (model.CandidateResume, error)
	GetResume(ctx context.Context, resumeID int) (model.CandidateResume, error)
	ListResumes(ctx context.Context) ([]model.CandidateResume, error)
	UpdateResume(ctx context.Context, resume model.CandidateResume) (model.CandidateResume, error)
	DeleteResume(ctx context.Context, resumeID int) error

	SaveJobAd(ctx context.Context, jobAd model.JobAd) (model.JobAd, error)
	GetJobAd(ctx context.Context, jobAdID int) (model.JobAd, error)
	ListJobAds(ctx context.Context) ([]model.JobAd, error)
	UpdateJobAd(ctx context.Context, jobAd model.JobAd) (model.JobAd, error)
	DeleteJobAd(ctx context.Context, jobAdID int) error
}

type LLMProvider interface {
	ParseResume(ctx context.Context, file model.File) (model.CandidateResume, error)
	ParseJobAd(ctx context.Context, file model.File) (model.JobAd, error)

	AdaptResume(ctx context.Context, jobAd model.JobAd, resumes []model.CandidateResume) (model.CandidateAdaptedResume, error)
}

type LLMProviderFactory interface {
	GetProvider(providerName string) (LLMProvider, error)
}

type SynthesizerService struct {
	llmFactory LLMProviderFactory
	repository ResumeRepository
}

func NewSynthesizerService(factory LLMProviderFactory, repo ResumeRepository) *SynthesizerService {
	return &SynthesizerService{
		llmFactory: factory,
		repository: repo,
	}
}

func (s *SynthesizerService) ParseResume(ctx context.Context, file model.File, providerName string) (model.CandidateResume, error) {
	provider, err := s.llmFactory.GetProvider(providerName)
	if err != nil {
		return model.CandidateResume{}, fmt.Errorf("could not get llm provider %s: %w", providerName, err)
	}
	resume, err := provider.ParseResume(ctx, file)
	if err != nil {
		return model.CandidateResume{}, fmt.Errorf("could not parse resume with llm provider %s: %w", providerName, err)
	}
	return s.repository.SaveResume(ctx, resume)
}

func (s *SynthesizerService) GetResume(ctx context.Context, resumeID int) (model.CandidateResume, error) {
	return s.repository.GetResume(ctx, resumeID)
}

func (s *SynthesizerService) UpdateResume(ctx context.Context, resume model.CandidateResume) (model.CandidateResume, error) {
	return s.repository.UpdateResume(ctx, resume)
}

func (s *SynthesizerService) DeleteResume(ctx context.Context, resumeID int) error {
	return s.repository.DeleteResume(ctx, resumeID)
}

func (s *SynthesizerService) ListResumes(ctx context.Context) ([]model.CandidateResume, error) {
	return s.repository.ListResumes(ctx)
}

// --- CRUD Operations for JobAds ---

func (s *SynthesizerService) ParseJobAd(ctx context.Context, file model.File, providerName string) (model.JobAd, error) {
	provider, err := s.llmFactory.GetProvider(providerName)
	if err != nil {
		return model.JobAd{}, fmt.Errorf("could not get llm provider %s: %w", providerName, err)
	}
	jobAd, err := provider.ParseJobAd(ctx, file)
	if err != nil {
		return model.JobAd{}, fmt.Errorf("could not parse job ad with llm provider %s: %w", providerName, err)
	}
	return s.repository.SaveJobAd(ctx, jobAd)
}

func (s *SynthesizerService) GetJobAd(ctx context.Context, jobID int) (model.JobAd, error) {
	return s.repository.GetJobAd(ctx, jobID)
}

func (s *SynthesizerService) UpdateJobAd(ctx context.Context, jobAd model.JobAd) (model.JobAd, error) {
	return s.repository.UpdateJobAd(ctx, jobAd)
}

func (s *SynthesizerService) DeleteJobAd(ctx context.Context, jobID int) error {
	return s.repository.DeleteJobAd(ctx, jobID)
}

func (s *SynthesizerService) ListJobAds(ctx context.Context) ([]model.JobAd, error) {
	return s.repository.ListJobAds(ctx)
}

// --- CRUD Operations for CandidateAdaptedResumes ---

func (s *SynthesizerService) AdaptResume(ctx context.Context, jobAdID int, resumeIDs []int, providerName string) (model.CandidateAdaptedResume, error) {
	jobAd, err := s.repository.GetJobAd(ctx, jobAdID)
	if err != nil {
		return model.CandidateAdaptedResume{}, fmt.Errorf("failed to retrieve job ad with ID %d: %w", jobAdID, err)
	}

	resumes := make([]model.CandidateResume, len(resumeIDs))
	for i, resumeID := range resumeIDs {
		resume, err := s.repository.GetResume(ctx, resumeID)
		if err != nil {
			return model.CandidateAdaptedResume{}, fmt.Errorf("failed to retrieve resume with ID %d: %w", resumeID, err)
		}
		resumes[i] = resume
	}

	if len(resumes) == 0 {
		return model.CandidateAdaptedResume{}, fmt.Errorf("at least one resume must be provided for adaptation")
	}

	provider, err := s.llmFactory.GetProvider(providerName)
	if err != nil {
		return model.CandidateAdaptedResume{}, fmt.Errorf("could not get LLM provider '%s': %w", providerName, err)
	}

	adapted, err := provider.AdaptResume(ctx, jobAd, resumes)
	if err != nil {
		return model.CandidateAdaptedResume{}, fmt.Errorf("LLM failed to adapt resume: %w", err)
	}

	return s.repository.SaveAdaptedResume(ctx, adapted)
}

func (s *SynthesizerService) GetAdaptedResume(ctx context.Context, adaptedResumeID int) (model.CandidateAdaptedResume, error) {
	return s.repository.GetAdaptedResume(ctx, adaptedResumeID)
}

func (s *SynthesizerService) UpdateAdaptedResume(ctx context.Context, adaptedResume model.CandidateAdaptedResume) (model.CandidateAdaptedResume, error) {
	return s.repository.UpdateAdaptedResume(ctx, adaptedResume)
}

func (s *SynthesizerService) DeleteAdaptedResume(ctx context.Context, adaptedResumeID int) error {
	return s.repository.DeleteAdaptedResume(ctx, adaptedResumeID)
}

func (s *SynthesizerService) ListAdaptedResumes(ctx context.Context) ([]model.CandidateAdaptedResume, error) {
	return s.repository.ListAdaptedResumes(ctx)
}
