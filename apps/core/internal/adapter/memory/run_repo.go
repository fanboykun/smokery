package memory

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

type RunRepo struct{ s *Store }

func NewRunRepo(s *Store) *RunRepo { return &RunRepo{s: s} }

var _ port.RunRepo = (*RunRepo)(nil)

func (r *RunRepo) Create(_ context.Context, projectID uuid.UUID, planID *uuid.UUID) (*model.Run, error) {
	run := model.Run{ID: uuid.New(), ProjectID: projectID, PlanID: planID, Status: "pending", CreatedAt: time.Now()}
	r.s.mu.Lock()
	r.s.runs[run.ID] = run
	r.s.mu.Unlock()
	return &run, nil
}

func (r *RunRepo) Get(_ context.Context, id uuid.UUID) (*model.Run, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	run, ok := r.s.runs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &run, nil
}

func (r *RunRepo) ListByProject(_ context.Context, projectID uuid.UUID) ([]model.Run, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	var out []model.Run
	for _, run := range r.s.runs {
		if run.ProjectID == projectID {
			out = append(out, run)
		}
	}
	return out, nil
}

func (r *RunRepo) UpdateStatus(_ context.Context, id uuid.UUID, status string, startedAt, finishedAt *time.Time) (*model.Run, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	run, ok := r.s.runs[id]
	if !ok {
		return nil, ErrNotFound
	}
	run.Status = status
	run.StartedAt = startedAt
	run.FinishedAt = finishedAt
	r.s.runs[id] = run
	return &run, nil
}

func (r *RunRepo) SaveResult(_ context.Context, runID uuid.UUID, result []byte) (*model.StoredRunResult, error) {
	sr := model.StoredRunResult{ID: uuid.New(), RunID: runID, Result: result, CreatedAt: time.Now()}
	r.s.mu.Lock()
	r.s.runResults[runID] = sr
	r.s.mu.Unlock()
	return &sr, nil
}

func (r *RunRepo) GetResult(_ context.Context, runID uuid.UUID) (*model.StoredRunResult, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	sr, ok := r.s.runResults[runID]
	if !ok {
		return nil, ErrNotFound
	}
	return &sr, nil
}

func (r *RunRepo) DeleteOlderThan(_ context.Context, before time.Time) (int, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	count := 0
	for id, run := range r.s.runs {
		if run.CreatedAt.Before(before) {
			delete(r.s.runs, id)
			delete(r.s.runResults, id)
			count++
		}
	}
	return count, nil
}
