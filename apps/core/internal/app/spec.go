package app

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
	"github.com/fanboykun/smokery/apps/core/internal/spec"
)

type SpecService struct {
	specs      port.SpecRepo
	operations port.OperationRepo
}

func NewSpecService(s port.SpecRepo, o port.OperationRepo) *SpecService {
	return &SpecService{specs: s, operations: o}
}

type ImportSpecResult struct {
	Spec     *model.Spec    `json:"spec"`
	Analysis *spec.Analysis `json:"analysis"`
}

func (s *SpecService) Import(ctx context.Context, projectID uuid.UUID, raw []byte) (*ImportSpecResult, error) {
	analysis, err := spec.Parse(raw)
	if err != nil {
		return nil, err
	}
	analysisJSON, _ := json.Marshal(analysis)
	created, err := s.specs.Create(ctx, model.Spec{
		ProjectID: projectID,
		Version:   analysis.Version,
		Title:     analysis.Title,
		Raw:       raw,
		Analysis:  analysisJSON,
	})
	if err != nil {
		return nil, err
	}
	for _, op := range analysis.Operations {
		_, _ = s.operations.Create(ctx, model.Operation{
			SpecID:         created.ID,
			OperationID:    op.OperationID,
			Method:         op.Method,
			Path:           op.Path,
			Summary:        op.Summary,
			Tags:           op.Tags,
			Classification: op.Classification,
			IsDestructive:  op.IsDestructive,
		})
	}
	return &ImportSpecResult{Spec: created, Analysis: analysis}, nil
}

func (s *SpecService) Get(ctx context.Context, id uuid.UUID) (*model.Spec, error) {
	return s.specs.Get(ctx, id)
}

func (s *SpecService) ListByProject(ctx context.Context, projectID uuid.UUID) ([]model.Spec, error) {
	return s.specs.ListByProject(ctx, projectID)
}
