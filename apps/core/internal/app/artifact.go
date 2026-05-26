package app

import (
	"context"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

type ArtifactService struct {
	artifacts port.ArtifactRepo
}

func NewArtifactService(a port.ArtifactRepo) *ArtifactService {
	return &ArtifactService{artifacts: a}
}

func (s *ArtifactService) ListByRun(ctx context.Context, runID uuid.UUID) ([]model.Artifact, error) {
	return s.artifacts.ListByRun(ctx, runID)
}
