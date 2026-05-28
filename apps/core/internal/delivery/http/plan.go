package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/app"
	"github.com/fanboykun/smokery/apps/core/internal/compiler"
	"github.com/fanboykun/smokery/apps/core/internal/model"
)

type PlanPreviewInput struct {
	ProjectIDParam
	Body model.ProjectConfig
}

type PlanPreviewOutput struct {
	Body compiler.Output
}

func RegisterPlan(api huma.API, svc *app.PlanService) {
	huma.Register(api, huma.Operation{
		OperationID: "preview-plan",
		Method:      http.MethodPost,
		Path:        "/api/projects/{project-id}/plan/preview",
		Summary:     "Preview compiled plan without persisting",
	}, func(ctx context.Context, in *PlanPreviewInput) (*PlanPreviewOutput, error) {
		projectID, err := uuid.Parse(in.ProjectID)
		if err != nil {
			return nil, ErrBadRequest("plan/preview", "invalid project id")
		}
		out, err := svc.Preview(ctx, projectID, in.Body)
		if err != nil {
			return nil, ErrInternal("plan/preview", "compilation failed", err)
		}
		return &PlanPreviewOutput{Body: *out}, nil
	})
}
