package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/app"
)

func RegisterFailureClassifications(api huma.API, svc *app.FailureClassificationService) {
	huma.Register(api, huma.Operation{
		OperationID: "classify-run-failure",
		Method:      http.MethodPut,
		Path:        "/api/runs/{id}/failure-classification",
		Summary:     "Classify a run failure",
	}, func(ctx context.Context, in *ClassifyRunInput) (*FailureClassificationOutput, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, ErrBadRequest("failure-classification", "invalid run id")
		}
		fc, err := svc.Classify(ctx, id, in.Body.Classification, in.Body.Assignee, in.Body.Note, in.Body.Author)
		if err != nil {
			return nil, ErrInternal("failure-classification", "failed to classify", err)
		}
		return &FailureClassificationOutput{Body: *fc}, nil
	})

	huma.Get(api, "/api/runs/{id}/failure-classification", func(ctx context.Context, in *IDParam) (*FailureClassificationOutput, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, ErrBadRequest("failure-classification", "invalid run id")
		}
		fc, err := svc.GetByRun(ctx, id)
		if err != nil {
			return nil, ErrInternal("failure-classification", "failed to get classification", err)
		}
		if fc == nil {
			return nil, ErrNotFound("failure-classification", "no classification found")
		}
		return &FailureClassificationOutput{Body: *fc}, nil
	})

	huma.Get(api, "/api/projects/{project-id}/failure-classifications", func(ctx context.Context, in *ProjectIDParam) (*FailureClassificationListOutput, error) {
		projectID, err := uuid.Parse(in.ProjectID)
		if err != nil {
			return nil, ErrBadRequest("failure-classifications", "invalid project id")
		}
		fcs, err := svc.ListByProject(ctx, projectID)
		if err != nil {
			return nil, ErrInternal("failure-classifications", "failed to list", err)
		}
		return &FailureClassificationListOutput{Body: fcs}, nil
	})
}
