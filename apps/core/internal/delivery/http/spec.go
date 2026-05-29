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

type ImportSpecFromURLInput struct {
	ProjectIDParam
	Body struct {
		URL     string            `json:"url" minLength:"1" doc:"URL to download the OpenAPI spec from"`
		Headers map[string]string `json:"headers,omitempty" doc:"Custom headers for downloading (e.g. Authorization)"`
	}
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

	// Import from raw body (paste)
	huma.Register(api, huma.Operation{
		OperationID: "import-spec",
		Method:      http.MethodPost,
		Path:        "/api/projects/{project-id}/specs",
		Summary:     "Import OpenAPI spec from raw body",
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

	// Import from URL
	huma.Register(api, huma.Operation{
		OperationID: "import-spec-from-url",
		Method:      http.MethodPost,
		Path:        "/api/projects/{project-id}/specs/from-url",
		Summary:     "Import OpenAPI spec from URL",
	}, func(ctx context.Context, in *ImportSpecFromURLInput) (*SpecAnalysisOutput, error) {
		projectID, err := uuid.Parse(in.ProjectID)
		if err != nil {
			return nil, ErrBadRequest("specs", "invalid project id")
		}
		res, err := svc.ImportFromURL(ctx, projectID, in.Body.URL, in.Body.Headers)
		if err != nil {
			return nil, ErrUnprocessable("specs", "failed to import spec from url", err)
		}
		return &SpecAnalysisOutput{Body: *res.Analysis}, nil
	})
}
