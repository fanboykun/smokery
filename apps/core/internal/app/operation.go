package app

import (
	"context"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

type OperationService struct {
	ops port.OperationRepo
}

func NewOperationService(o port.OperationRepo) *OperationService {
	return &OperationService{ops: o}
}

func (s *OperationService) ListBySpec(ctx context.Context, specID uuid.UUID) ([]model.Operation, error) {
	return s.ops.ListBySpec(ctx, specID)
}

func (s *OperationService) UpdateClassification(ctx context.Context, id uuid.UUID, classification string, isDestructive bool) (*model.Operation, error) {
	return s.ops.UpdateClassification(ctx, id, classification, isDestructive)
}

func (s *OperationService) UpdateOverrides(ctx context.Context, id uuid.UUID, overrides []byte) (*model.Operation, error) {
	return s.ops.UpdateOverrides(ctx, id, overrides)
}
