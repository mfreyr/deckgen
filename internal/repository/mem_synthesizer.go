package storage

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/mfreyr/deckgen/internal/model"
)

// MemoryResumeRepo is an in-memory implementation of the ResumeRepository interface.
// It uses maps to store data and is safe for concurrent use.
type MemoryResumeRepo struct {
	mu sync.RWMutex

	resumes        map[int]model.CandidateResume
	jobAds         map[int]model.JobAd
	adaptedResumes map[int]model.CandidateAdaptedResume

	nextResumeID        int
	nextJobAdID         int
	nextAdaptedResumeID int
}

// NewMemoryResumeRepo creates and initializes a new in-memory repository.
func NewMemoryResumeRepo() *MemoryResumeRepo {
	return &MemoryResumeRepo{
		resumes:        make(map[int]model.CandidateResume),
		jobAds:         make(map[int]model.JobAd),
		adaptedResumes: make(map[int]model.CandidateAdaptedResume),

		nextResumeID:        1,
		nextJobAdID:         1,
		nextAdaptedResumeID: 1,
	}
}

// --- AdaptedResume Methods ---

func (r *MemoryResumeRepo) SaveAdaptedResume(ctx context.Context, adaptedResume model.CandidateAdaptedResume) (model.CandidateAdaptedResume, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	adaptedResume.ID = r.nextAdaptedResumeID
	r.adaptedResumes[adaptedResume.ID] = adaptedResume
	r.nextAdaptedResumeID++

	return adaptedResume, nil
}

func (r *MemoryResumeRepo) GetAdaptedResume(ctx context.Context, adaptedResumeID int) (model.CandidateAdaptedResume, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resume, ok := r.adaptedResumes[adaptedResumeID]
	if !ok {
		return model.CandidateAdaptedResume{}, fmt.Errorf("adapted resume with ID %d not found", adaptedResumeID)
	}
	return resume, nil
}

func (r *MemoryResumeRepo) ListAdaptedResumes(ctx context.Context) ([]model.CandidateAdaptedResume, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resumes := make([]model.CandidateAdaptedResume, 0, len(r.adaptedResumes))
	for _, resume := range r.adaptedResumes {
		resumes = append(resumes, resume)
	}
	// Sort for consistent ordering
	sort.Slice(resumes, func(i, j int) bool {
		return resumes[i].ID < resumes[j].ID
	})
	return resumes, nil
}

func (r *MemoryResumeRepo) UpdateAdaptedResume(ctx context.Context, adaptedResume model.CandidateAdaptedResume) (model.CandidateAdaptedResume, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.adaptedResumes[adaptedResume.ID]; !ok {
		return model.CandidateAdaptedResume{}, fmt.Errorf("adapted resume with ID %d not found for update", adaptedResume.ID)
	}
	r.adaptedResumes[adaptedResume.ID] = adaptedResume
	return adaptedResume, nil
}

func (r *MemoryResumeRepo) DeleteAdaptedResume(ctx context.Context, adaptedResumeID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.adaptedResumes[adaptedResumeID]; !ok {
		return fmt.Errorf("adapted resume with ID %d not found for deletion", adaptedResumeID)
	}
	delete(r.adaptedResumes, adaptedResumeID)
	return nil
}

// --- CandidateResume Methods ---

func (r *MemoryResumeRepo) SaveResume(ctx context.Context, resume model.CandidateResume) (model.CandidateResume, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	resume.ID = r.nextResumeID
	r.resumes[resume.ID] = resume
	r.nextResumeID++

	return resume, nil
}

func (r *MemoryResumeRepo) GetResume(ctx context.Context, resumeID int) (model.CandidateResume, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resume, ok := r.resumes[resumeID]
	if !ok {
		return model.CandidateResume{}, fmt.Errorf("resume with ID %d not found", resumeID)
	}
	return resume, nil
}

func (r *MemoryResumeRepo) ListResumes(ctx context.Context) ([]model.CandidateResume, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resumes := make([]model.CandidateResume, 0, len(r.resumes))
	for _, resume := range r.resumes {
		resumes = append(resumes, resume)
	}
	sort.Slice(resumes, func(i, j int) bool {
		return resumes[i].ID < resumes[j].ID
	})
	return resumes, nil
}

func (r *MemoryResumeRepo) UpdateResume(ctx context.Context, resume model.CandidateResume) (model.CandidateResume, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.resumes[resume.ID]; !ok {
		return model.CandidateResume{}, fmt.Errorf("resume with ID %d not found for update", resume.ID)
	}
	r.resumes[resume.ID] = resume
	return resume, nil
}

func (r *MemoryResumeRepo) DeleteResume(ctx context.Context, resumeID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.resumes[resumeID]; !ok {
		return fmt.Errorf("resume with ID %d not found for deletion", resumeID)
	}
	delete(r.resumes, resumeID)
	return nil
}

// --- JobAd Methods ---

func (r *MemoryResumeRepo) SaveJobAd(ctx context.Context, jobAd model.JobAd) (model.JobAd, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	jobAd.ID = r.nextJobAdID
	r.jobAds[jobAd.ID] = jobAd
	r.nextJobAdID++

	return jobAd, nil
}

func (r *MemoryResumeRepo) GetJobAd(ctx context.Context, jobAdID int) (model.JobAd, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	jobAd, ok := r.jobAds[jobAdID]
	if !ok {
		return model.JobAd{}, fmt.Errorf("job ad with ID %d not found", jobAdID)
	}
	return jobAd, nil
}

func (r *MemoryResumeRepo) ListJobAds(ctx context.Context) ([]model.JobAd, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	jobAds := make([]model.JobAd, 0, len(r.jobAds))
	for _, jobAd := range r.jobAds {
		jobAds = append(jobAds, jobAd)
	}
	sort.Slice(jobAds, func(i, j int) bool {
		return jobAds[i].ID < jobAds[j].ID
	})
	return jobAds, nil
}

func (r *MemoryResumeRepo) UpdateJobAd(ctx context.Context, jobAd model.JobAd) (model.JobAd, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.jobAds[jobAd.ID]; !ok {
		return model.JobAd{}, fmt.Errorf("job ad with ID %d not found for update", jobAd.ID)
	}
	r.jobAds[jobAd.ID] = jobAd
	return jobAd, nil
}

func (r *MemoryResumeRepo) DeleteJobAd(ctx context.Context, jobAdID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.jobAds[jobAdID]; !ok {
		return fmt.Errorf("job ad with ID %d not found for deletion", jobAdID)
	}
	delete(r.jobAds, jobAdID)
	return nil
}
