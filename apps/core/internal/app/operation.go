package app

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
	"github.com/fanboykun/smokery/apps/core/internal/spec"
)

type OperationService struct {
	specs port.SpecRepo
	ops   port.OperationRepo
}

func NewOperationService(specs port.SpecRepo, ops port.OperationRepo) *OperationService {
	return &OperationService{specs: specs, ops: ops}
}

func (s *OperationService) ListBySpec(ctx context.Context, specID uuid.UUID) ([]model.Operation, error) {
	return s.ops.ListBySpec(ctx, specID)
}

func (s *OperationService) CanvasBySpec(ctx context.Context, specID uuid.UUID) ([]spec.OperationInfo, error) {
	sp, err := s.specs.Get(ctx, specID)
	if err != nil {
		return nil, err
	}
	var analysis spec.Analysis
	if err := json.Unmarshal(sp.Analysis, &analysis); err != nil {
		return nil, err
	}
	return analysis.Operations, nil
}

func (s *OperationService) UpdateClassification(ctx context.Context, id uuid.UUID, classification string, isDestructive bool) (*model.Operation, error) {
	return s.ops.UpdateClassification(ctx, id, classification, isDestructive)
}

func (s *OperationService) UpdateOverrides(ctx context.Context, id uuid.UUID, overrides []byte) (*model.Operation, error) {
	return s.ops.UpdateOverrides(ctx, id, overrides)
}
