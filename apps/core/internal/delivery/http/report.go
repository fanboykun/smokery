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
			return nil, ErrBadRequest("reports", "invalid id")
		}
		v, err := svc.DebugView(ctx, id)
		if err != nil {
			return nil, ErrNotFound("reports", "result not found")
		}
		return &DebugReportOutput{Body: *v}, nil
	})

	huma.Get(api, "/api/runs/{id}/report/ci", func(ctx context.Context, in *IDParam) (*CIReportOutput, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, ErrBadRequest("reports", "invalid id")
		}
		v, err := svc.CISummary(ctx, id)
		if err != nil {
			return nil, ErrNotFound("reports", "result not found")
		}
		return &CIReportOutput{Body: *v}, nil
	})

	huma.Get(api, "/api/runs/{id}/report/mermaid", func(ctx context.Context, in *IDParam) (*MermaidOutput, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, ErrBadRequest("reports", "invalid id")
		}
		s, err := svc.Mermaid(ctx, id)
		if err != nil {
			return nil, ErrNotFound("reports", "result not found")
		}
		return &MermaidOutput{Body: s}, nil
	})

	huma.Get(api, "/api/runs/{id}/report/contract", func(ctx context.Context, in *IDParam) (*ContractReportOutput, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, ErrBadRequest("reports", "invalid id")
		}
		v, err := svc.ContractView(ctx, id)
		if err != nil {
			return nil, ErrNotFound("reports", "result not found")
		}
		return &ContractReportOutput{Body: *v}, nil
	})

	huma.Get(api, "/api/runs/{id}/report/analyst", func(ctx context.Context, in *IDParam) (*AnalystReportOutput, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, ErrBadRequest("reports", "invalid id")
		}
		v, err := svc.AnalystView(ctx, id)
		if err != nil {
			return nil, ErrNotFound("reports", "result not found")
		}
		return &AnalystReportOutput{Body: *v}, nil
	})

	huma.Get(api, "/api/runs/{id}/report/qa", func(ctx context.Context, in *IDParam) (*QAReportOutput, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, ErrBadRequest("reports", "invalid id")
		}
		v, err := svc.QAView(ctx, id)
		if err != nil {
			return nil, ErrNotFound("reports", "result not found")
		}
		return &QAReportOutput{Body: *v}, nil
	})
}
