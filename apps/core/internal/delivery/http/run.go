package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/app"
)

// RegisterRuns registers run lifecycle endpoints.
func RegisterRuns(api huma.API, svc *app.RunService) {
	huma.Register(api, huma.Operation{
		OperationID: "create-run",
		Method:      http.MethodPost,
		Path:        "/api/projects/{project-id}/runs",
		Summary:     "Create and start a run",
	}, func(ctx context.Context, in *CreateRunInput) (*RunOutput, error) {
		projectID, err := uuid.Parse(in.ProjectID)
		if err != nil {
			return nil, ErrBadRequest("runs", "invalid project id")
		}
		var planID *uuid.UUID
		if in.Body.PlanID != "" {
			pid, err := uuid.Parse(in.Body.PlanID)
			if err == nil {
				planID = &pid
			}
		}
		run, err := svc.Start(ctx, app.StartRunInput{
			ProjectID: projectID,
			PlanID:    planID,
			Plan:      in.Body.Plan,
		})
		if err != nil {
			return nil, ErrInternal("runs", "failed to start run", err)
		}
		return &RunOutput{Body: *run}, nil
	})

	huma.Get(api, "/api/projects/{project-id}/runs", func(ctx context.Context, in *ProjectIDParam) (*RunListOutput, error) {
		projectID, err := uuid.Parse(in.ProjectID)
		if err != nil {
			return nil, ErrBadRequest("runs", "invalid project id")
		}
		runs, err := svc.ListByProject(ctx, projectID)
		if err != nil {
			return nil, ErrInternal("runs", "failed to list runs", err)
		}
		return &RunListOutput{Body: runs}, nil
	})

	huma.Get(api, "/api/runs/{id}", func(ctx context.Context, in *IDParam) (*RunOutput, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, ErrBadRequest("runs/{id}", "invalid run id")
		}
		run, err := svc.Get(ctx, id)
		if err != nil {
			return nil, ErrNotFound("runs/{id}", "run not found")
		}
		return &RunOutput{Body: *run}, nil
	})

	huma.Get(api, "/api/runs/{id}/result", func(ctx context.Context, in *IDParam) (*RunResultOutput, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, ErrBadRequest("runs/{id}/result", "invalid run id")
		}
		res, err := svc.GetResult(ctx, id)
		if err != nil {
			return nil, ErrNotFound("runs/{id}/result", "result not found")
		}
		return &RunResultOutput{Body: *res}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "cancel-run",
		Method:      http.MethodPost,
		Path:        "/api/runs/{id}/cancel",
		Summary:     "Cancel a running run",
	}, func(ctx context.Context, in *IDParam) (*RunOutput, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, ErrBadRequest("runs/{id}/cancel", "invalid run id")
		}
		run, err := svc.Cancel(ctx, id)
		if err != nil {
			return nil, ErrInternal("runs/{id}/cancel", "failed to cancel run", err)
		}
		return &RunOutput{Body: *run}, nil
	})
}
