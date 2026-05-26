package app

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
	"github.com/fanboykun/smokery/apps/core/internal/report"
)

type ReportService struct {
	runs port.RunRepo
}

func NewReportService(r port.RunRepo) *ReportService {
	return &ReportService{runs: r}
}

func (s *ReportService) loadResult(ctx context.Context, runID uuid.UUID) (*model.RunResult, error) {
	stored, err := s.runs.GetResult(ctx, runID)
	if err != nil {
		return nil, err
	}
	if stored == nil || len(stored.Result) == 0 {
		return nil, errors.New("no result")
	}
	var rr model.RunResult
	if err := json.Unmarshal(stored.Result, &rr); err != nil {
		return nil, err
	}
	return &rr, nil
}

func (s *ReportService) DebugView(ctx context.Context, runID uuid.UUID) (*report.DebugView, error) {
	rr, err := s.loadResult(ctx, runID)
	if err != nil {
		return nil, err
	}
	return report.GenerateDebugView(rr), nil
}

func (s *ReportService) CISummary(ctx context.Context, runID uuid.UUID) (*report.CISummary, error) {
	rr, err := s.loadResult(ctx, runID)
	if err != nil {
		return nil, err
	}
	return report.GenerateCISummary(rr), nil
}

func (s *ReportService) Mermaid(ctx context.Context, runID uuid.UUID) (string, error) {
	rr, err := s.loadResult(ctx, runID)
	if err != nil {
		return "", err
	}
	return report.GenerateMermaidDiagram(rr), nil
}
