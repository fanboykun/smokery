package app

import (
	"context"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

type RunService struct {
	runs port.RunRepo
	jobs port.JobEnqueuer
}

func NewRunService(r port.RunRepo, j port.JobEnqueuer) *RunService {
	return &RunService{runs: r, jobs: j}
}

type StartRunInput struct {
	ProjectID uuid.UUID
	PlanID    *uuid.UUID
	Plan      *model.SmokePlan
}

func (s *RunService) Start(ctx context.Context, in StartRunInput) (*model.Run, error) {
	run, err := s.runs.Create(ctx, in.ProjectID, in.PlanID)
	if err != nil {
		return nil, err
	}
	if err := s.jobs.EnqueueRun(context.Background(), run.ID, in.Plan); err != nil {
		return nil, err
	}
	return run, nil
}

func (s *RunService) Get(ctx context.Context, id uuid.UUID) (*model.Run, error) {
	return s.runs.Get(ctx, id)
}

func (s *RunService) ListByProject(ctx context.Context, projectID uuid.UUID) ([]model.Run, error) {
	return s.runs.ListByProject(ctx, projectID)
}

func (s *RunService) GetResult(ctx context.Context, runID uuid.UUID) (*model.StoredRunResult, error) {
	return s.runs.GetResult(ctx, runID)
}
