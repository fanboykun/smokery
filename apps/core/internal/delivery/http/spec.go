package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/app"
	"github.com/fanboykun/smokery/apps/core/internal/model"
)

type SpecListOutput struct {
	Body []model.Spec
}

// RegisterSpecs registers spec import and list endpoints.
func RegisterSpecs(api huma.API, svc *app.SpecService) {
	huma.Get(api, "/api/projects/{project-id}/specs", func(ctx context.Context, in *ProjectIDParam) (*SpecListOutput, error) {
		projectID, err := uuid.Parse(in.ProjectID)
		if err != nil {
			return nil, ErrBadRequest("specs", "invalid project id")
		}
		specs, err := svc.ListByProject(ctx, projectID)
		if err != nil {
			return nil, ErrInternal("specs", "failed to list specs", err)
		}
		return &SpecListOutput{Body: specs}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "import-spec",
		Method:      http.MethodPost,
		Path:        "/api/projects/{project-id}/specs",
		Summary:     "Import OpenAPI spec",
	}, func(ctx context.Context, in *ImportSpecInput) (*SpecAnalysisOutput, error) {
		projectID, err := uuid.Parse(in.ProjectID)
		if err != nil {
			return nil, ErrBadRequest("specs", "invalid project id")
		}
		res, err := svc.Import(ctx, projectID, in.RawBody)
		if err != nil {
			return nil, ErrUnprocessable("specs", "failed to import spec", err)
		}
		return &SpecAnalysisOutput{Body: *res.Analysis}, nil
	})
}
