package http

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/app"
)

// RegisterReports registers run report views.
func RegisterReports(api huma.API, svc *app.ReportService) {
	huma.Get(api, "/api/runs/{id}/report/debug", func(ctx context.Context, in *IDParam) (*DebugReportOutput, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id")
		}
		v, err := svc.DebugView(ctx, id)
		if err != nil {
			return nil, huma.Error404NotFound("result not found")
		}
		return &DebugReportOutput{Body: *v}, nil
	})

	huma.Get(api, "/api/runs/{id}/report/ci", func(ctx context.Context, in *IDParam) (*CIReportOutput, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id")
		}
		v, err := svc.CISummary(ctx, id)
		if err != nil {
			return nil, huma.Error404NotFound("result not found")
		}
		return &CIReportOutput{Body: *v}, nil
	})

	huma.Get(api, "/api/runs/{id}/report/mermaid", func(ctx context.Context, in *IDParam) (*MermaidOutput, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id")
		}
		s, err := svc.Mermaid(ctx, id)
		if err != nil {
			return nil, huma.Error404NotFound("result not found")
		}
		return &MermaidOutput{Body: s}, nil
	})
}
