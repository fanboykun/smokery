package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

func (s *SpecService) ImportFromURL(ctx context.Context, projectID uuid.UUID, url string, headers map[string]string) (*ImportSpecResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download spec: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download spec: status %d", resp.StatusCode)
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 50<<20)) // 50MB limit
	if err != nil {
		return nil, fmt.Errorf("read spec body: %w", err)
	}
	return s.Import(ctx, projectID, raw)
}
