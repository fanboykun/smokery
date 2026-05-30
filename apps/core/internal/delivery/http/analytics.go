package http

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/app"
)

type AnalyticsRangeInput struct {
	ProjectIDParam
	Range string `query:"range" default:"30d" doc:"Time range (7d, 30d, 90d)"`
}

type LatencyOutput struct {
	Body app.LatencyAnalytics
}

type FlakyOutput struct {
	Body app.FlakyAnalytics
}

type HealthTrendsOutput struct {
	Body app.HealthTrends
}

func RegisterAnalytics(api huma.API, svc *app.AnalyticsService) {
	huma.Get(api, "/api/projects/{project-id}/analytics/latency", func(ctx context.Context, in *AnalyticsRangeInput) (*LatencyOutput, error) {
		projectID, err := uuid.Parse(in.ProjectID)
		if err != nil {
			return nil, ErrBadRequest("analytics", "invalid project id")
		}
		data, err := svc.Latency(ctx, projectID, in.Range)
		if err != nil {
			return nil, ErrInternal("analytics", "failed to compute latency", err)
		}
		return &LatencyOutput{Body: *data}, nil
	})

	huma.Get(api, "/api/projects/{project-id}/analytics/flaky-operations", func(ctx context.Context, in *AnalyticsRangeInput) (*FlakyOutput, error) {
		projectID, err := uuid.Parse(in.ProjectID)
		if err != nil {
			return nil, ErrBadRequest("analytics", "invalid project id")
		}
		data, err := svc.FlakyOperations(ctx, projectID, in.Range)
		if err != nil {
			return nil, ErrInternal("analytics", "failed to compute flaky ops", err)
		}
		return &FlakyOutput{Body: *data}, nil
	})

	huma.Get(api, "/api/projects/{project-id}/analytics/health-trends", func(ctx context.Context, in *AnalyticsRangeInput) (*HealthTrendsOutput, error) {
		projectID, err := uuid.Parse(in.ProjectID)
		if err != nil {
			return nil, ErrBadRequest("analytics", "invalid project id")
		}
		data, err := svc.HealthTrends(ctx, projectID, in.Range)
		if err != nil {
			return nil, ErrInternal("analytics", "failed to compute health trends", err)
		}
		return &HealthTrendsOutput{Body: *data}, nil
	})
}
