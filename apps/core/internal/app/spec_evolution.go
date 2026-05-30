package app

import (
	"context"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

type SpecEvolutionService struct {
	specs      port.SpecRepo
	operations port.OperationRepo
}

func NewSpecEvolutionService(s port.SpecRepo, o port.OperationRepo) *SpecEvolutionService {
	return &SpecEvolutionService{specs: s, operations: o}
}

type SpecDiffChange struct {
	Type        string `json:"type"` // added, removed, modified
	OperationID string `json:"operation_id"`
	Path        string `json:"path"`
	Method      string `json:"method"`
	Breaking    bool   `json:"breaking"`
	Details     string `json:"details"`
}

type SpecDiff struct {
	FromSpecID string           `json:"from_spec_id"`
	ToSpecID   string           `json:"to_spec_id"`
	Changes    []SpecDiffChange `json:"changes"`
}

func (s *SpecEvolutionService) Diff(ctx context.Context, fromID, toID uuid.UUID) (*SpecDiff, error) {
	fromOps, err := s.operations.ListBySpec(ctx, fromID)
	if err != nil {
		return nil, err
	}
	toOps, err := s.operations.ListBySpec(ctx, toID)
	if err != nil {
		return nil, err
	}

	fromMap := map[string]model.Operation{}
	for _, op := range fromOps {
		fromMap[op.OperationID] = op
	}
	toMap := map[string]model.Operation{}
	for _, op := range toOps {
		toMap[op.OperationID] = op
	}

	diff := &SpecDiff{FromSpecID: fromID.String(), ToSpecID: toID.String()}

	// Added
	for _, op := range toOps {
		if _, exists := fromMap[op.OperationID]; !exists {
			diff.Changes = append(diff.Changes, SpecDiffChange{
				Type: "added", OperationID: op.OperationID,
				Path: op.Path, Method: op.Method, Breaking: false,
				Details: "New endpoint added",
			})
		}
	}
	// Removed (breaking)
	for _, op := range fromOps {
		if _, exists := toMap[op.OperationID]; !exists {
			diff.Changes = append(diff.Changes, SpecDiffChange{
				Type: "removed", OperationID: op.OperationID,
				Path: op.Path, Method: op.Method, Breaking: true,
				Details: "Endpoint removed",
			})
		}
	}
	// Modified
	for _, toOp := range toOps {
		if fromOp, exists := fromMap[toOp.OperationID]; exists {
			if fromOp.Path != toOp.Path || fromOp.Method != toOp.Method {
				diff.Changes = append(diff.Changes, SpecDiffChange{
					Type: "modified", OperationID: toOp.OperationID,
					Path: toOp.Path, Method: toOp.Method, Breaking: true,
					Details: "Path or method changed",
				})
			}
		}
	}
	return diff, nil
}

type AffectedFlow struct {
	FlowID        string `json:"flow_id"`
	FlowName      string `json:"flow_name"`
	AffectedSteps int    `json:"affected_steps"`
	Impact        string `json:"impact"`
}

type AffectedSuite struct {
	SuiteID            string `json:"suite_id"`
	SuiteName          string `json:"suite_name"`
	AffectedOperations int    `json:"affected_operations"`
	Impact             string `json:"impact"`
}

type ImpactAnalysis struct {
	SpecVersionID  string          `json:"spec_version_id"`
	AffectedFlows  []AffectedFlow  `json:"affected_flows"`
	AffectedSuites []AffectedSuite `json:"affected_suites"`
	AffectedRuns   int             `json:"affected_runs"`
	RiskAssessment string          `json:"risk_assessment"`
}

func (s *SpecEvolutionService) Impact(ctx context.Context, projectID, specID uuid.UUID) (*ImpactAnalysis, error) {
	// Get all specs for the project to find the previous one
	specs, err := s.specs.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var prevSpecID *uuid.UUID
	for i, sp := range specs {
		if sp.ID == specID && i > 0 {
			prevSpecID = &specs[i-1].ID
			break
		}
	}

	analysis := &ImpactAnalysis{SpecVersionID: specID.String()}

	if prevSpecID == nil {
		analysis.RiskAssessment = "low"
		return analysis, nil
	}

	diff, err := s.Diff(ctx, *prevSpecID, specID)
	if err != nil {
		return analysis, nil
	}

	breakingCount := 0
	for _, c := range diff.Changes {
		if c.Breaking {
			breakingCount++
		}
	}

	switch {
	case breakingCount >= 3:
		analysis.RiskAssessment = "high"
	case breakingCount >= 1:
		analysis.RiskAssessment = "medium"
	default:
		analysis.RiskAssessment = "low"
	}

	analysis.AffectedRuns = breakingCount * 10 // estimate
	return analysis, nil
}
