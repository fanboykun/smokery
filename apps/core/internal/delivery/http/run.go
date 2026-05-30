package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/app"
	"github.com/fanboykun/smokery/apps/core/internal/model"
)

const runEventTTL = time.Hour

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

	huma.Get(api, "/api/runs/{id}", func(ctx context.Context, in *RunMetaInput) (*RunMetaOutput, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, ErrBadRequest("runs/{id}", "invalid run id")
		}
		run, err := svc.Get(ctx, id)
		if err != nil {
			return nil, ErrNotFound("runs/{id}", "run not found")
		}

		var planID *string
		if run.PlanID != nil {
			value := run.PlanID.String()
			planID = &value
		}

		return &RunMetaOutput{Body: RunMeta{
			ID:                     run.ID.String(),
			ProjectID:              run.ProjectID.String(),
			PlanID:                 planID,
			Status:                 run.Status,
			CreatedAt:              run.CreatedAt,
			StartedAt:              run.StartedAt,
			FinishedAt:             run.FinishedAt,
			WebSocketURL:           runWebSocketURL(in, run.ID),
			FallbackPollIntervalMs: 1000,
			ExpiresAt:              run.CreatedAt.Add(runEventTTL),
		}}, nil
	})

	huma.Get(api, "/api/runs/{id}/result", func(ctx context.Context, in *IDParam) (*RunResultDetailsOutput, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, ErrBadRequest("runs/{id}/result", "invalid run id")
		}
		res, err := svc.GetResult(ctx, id)
		if err != nil {
			return nil, ErrNotFound("runs/{id}/result", "result not found")
		}

		var result model.RunResult
		if err := json.Unmarshal(res.Result, &result); err != nil {
			return nil, ErrInternal("runs/{id}/result", "failed to decode result", err)
		}
		return &RunResultDetailsOutput{Body: result}, nil
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

func runWebSocketURL(in *RunMetaInput, runID uuid.UUID) string {
	host := strings.TrimSpace(in.ForwardedHost)
	if host == "" {
		host = strings.TrimSpace(in.Host)
	}
	if host == "" {
		host = "localhost:8080"
	}

	proto := strings.TrimSpace(in.ForwardedProto)
	if proto == "" {
		proto = strings.TrimSpace(in.ForwardedScheme)
	}
	if proto == "" {
		proto = "http"
	}

	scheme := "ws"
	if strings.EqualFold(proto, "https") || strings.EqualFold(proto, "wss") {
		scheme = "wss"
	}

	return fmt.Sprintf("%s://%s/api/runs/%s/events", scheme, host, runID.String())
}
