package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/app"
)

// RegisterSpecs registers spec import endpoint.
func RegisterSpecs(api huma.API, svc *app.SpecService) {
	huma.Register(api, huma.Operation{
		OperationID: "import-spec",
		Method:      http.MethodPost,
		Path:        "/api/projects/{project-id}/specs",
		Summary:     "Import OpenAPI spec",
	}, func(ctx context.Context, in *ImportSpecInput) (*SpecAnalysisOutput, error) {
		projectID, err := uuid.Parse(in.ProjectID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid project id")
		}
		res, err := svc.Import(ctx, projectID, in.RawBody)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity("failed to import spec", err)
		}
		return &SpecAnalysisOutput{Body: *res.Analysis}, nil
	})
}
