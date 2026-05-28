package app

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/compiler"
	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
	"github.com/fanboykun/smokery/apps/core/internal/spec"
)

type PlanService struct {
	specs port.SpecRepo
	ops   port.OperationRepo
}

func NewPlanService(s port.SpecRepo, o port.OperationRepo) *PlanService {
	return &PlanService{specs: s, ops: o}
}

// Preview compiles a project config into a SmokePlan without persisting.
// It fetches operations from the project's latest spec.
func (s *PlanService) Preview(ctx context.Context, projectID uuid.UUID, cfg model.ProjectConfig) (*compiler.Output, error) {
	specs, err := s.specs.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if len(specs) == 0 {
		out := compiler.Compile(compiler.Input{Config: cfg})
		return &out, nil
	}
	// Use latest spec's analysis to get operations
	latest := specs[len(specs)-1]
	var analysis spec.Analysis
	if err := json.Unmarshal(latest.Analysis, &analysis); err != nil {
		out := compiler.Compile(compiler.Input{Config: cfg})
		return &out, nil
	}
	out := compiler.Compile(compiler.Input{Config: cfg, Operations: analysis.Operations})
	return &out, nil
}
