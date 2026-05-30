package app

import (
	"context"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

type FailureClassificationService struct {
	repo port.FailureClassificationRepo
}

func NewFailureClassificationService(r port.FailureClassificationRepo) *FailureClassificationService {
	return &FailureClassificationService{repo: r}
}

func (s *FailureClassificationService) Classify(ctx context.Context, runID uuid.UUID, classification, assignee, note, author string) (*model.FailureClassification, error) {
	return s.repo.Upsert(ctx, runID, classification, assignee, note, author)
}

func (s *FailureClassificationService) GetByRun(ctx context.Context, runID uuid.UUID) (*model.FailureClassification, error) {
	return s.repo.GetByRun(ctx, runID)
}

func (s *FailureClassificationService) ListByProject(ctx context.Context, projectID uuid.UUID) ([]model.FailureClassification, error) {
	return s.repo.ListByProject(ctx, projectID)
}
