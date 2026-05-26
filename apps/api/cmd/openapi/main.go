package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"

	"github.com/fanboykun/smokery/apps/api/internal/db"
	"github.com/fanboykun/smokery/apps/api/internal/model"
	"github.com/fanboykun/smokery/apps/api/internal/report"
	"github.com/fanboykun/smokery/apps/api/internal/spec"
)

// Reuse the same Input/Output types from the server
type CreateProjectInput struct {
	Body struct {
		Name        string `json:"name" minLength:"1" doc:"Project name"`
		Description string `json:"description" doc:"Project description"`
	}
}
type ProjectOutput struct{ Body db.Project }
type ProjectListOutput struct{ Body []db.Project }
type IDParam struct {
	ID string `path:"id" format:"uuid" doc:"Resource UUID"`
}
type ProjectIDParam struct {
	ProjectID string `path:"project-id" format:"uuid" doc:"Project UUID"`
}
type SpecIDParam struct {
	SpecID string `path:"spec-id" format:"uuid" doc:"Spec UUID"`
}
type UpdateProjectInput struct {
	IDParam
	Body struct {
		Name        string `json:"name" minLength:"1" doc:"Project name"`
		Description string `json:"description" doc:"Project description"`
	}
}
type ImportSpecInput struct {
	ProjectIDParam
	RawBody []byte
}
type SpecAnalysisOutput struct{ Body spec.Analysis }
type OperationListOutput struct{ Body []db.Operation }
type UpdateClassificationInput struct {
	IDParam
	Body struct {
		Classification string `json:"classification" minLength:"1" doc:"Operation classification"`
		IsDestructive  bool   `json:"is_destructive" doc:"Whether operation is destructive"`
	}
}
type OperationOutput struct{ Body db.Operation }
type UpdateOverridesInput struct {
	IDParam
	RawBody []byte
}
type CreateRunInput struct {
	ProjectIDParam
	Body struct {
		PlanID string           `json:"plan_id" doc:"Plan UUID"`
		Plan   *model.SmokePlan `json:"plan" doc:"Compiled smoke plan to execute"`
	}
}
type RunOutput struct{ Body db.Run }
type RunListOutput struct{ Body []db.Run }
type RunResultOutput struct{ Body db.RunResult }
type DebugReportOutput struct{ Body report.DebugView }
type CIReportOutput struct{ Body report.CISummary }
type MermaidOutput struct{ Body string }
type ArtifactListOutput struct{ Body []db.Artifact }
type CreateCommentInput struct {
	IDParam
	Body struct {
		Author string `json:"author" minLength:"1" doc:"Comment author"`
		Body   string `json:"body" minLength:"1" doc:"Comment body"`
	}
}
type CommentOutput struct{ Body db.Comment }
type CommentListOutput struct{ Body []db.Comment }
type HealthOutput struct {
	Body struct {
		Status string `json:"status" example:"ok" doc:"Service health status"`
	}
}

func main() {
	e := echo.New()
	cfg := huma.DefaultConfig("Smokery API", "1.0.0")
	cfg.Servers = []*huma.Server{{URL: "http://localhost:8080"}}
	api := humaecho.New(e, cfg)

	noop := func(ctx context.Context, input *struct{}) (*struct{}, error) { return nil, nil }
	_ = noop

	huma.Get(api, "/health", func(ctx context.Context, input *struct{}) (*HealthOutput, error) { return nil, nil })
	huma.Post(api, "/api/projects", func(ctx context.Context, input *CreateProjectInput) (*ProjectOutput, error) { return nil, nil })
	huma.Get(api, "/api/projects", func(ctx context.Context, input *struct{}) (*ProjectListOutput, error) { return nil, nil })
	huma.Get(api, "/api/projects/{id}", func(ctx context.Context, input *IDParam) (*ProjectOutput, error) { return nil, nil })
	huma.Put(api, "/api/projects/{id}", func(ctx context.Context, input *UpdateProjectInput) (*ProjectOutput, error) { return nil, nil })
	huma.Delete(api, "/api/projects/{id}", func(ctx context.Context, input *IDParam) (*struct{}, error) { return nil, nil })

	huma.Register(api, huma.Operation{OperationID: "import-spec", Method: "POST", Path: "/api/projects/{project-id}/specs", Summary: "Import OpenAPI spec"},
		func(ctx context.Context, input *ImportSpecInput) (*SpecAnalysisOutput, error) { return nil, nil })

	huma.Get(api, "/api/specs/{spec-id}/operations", func(ctx context.Context, input *SpecIDParam) (*OperationListOutput, error) { return nil, nil })
	huma.Put(api, "/api/operations/{id}/classification", func(ctx context.Context, input *UpdateClassificationInput) (*OperationOutput, error) { return nil, nil })
	huma.Register(api, huma.Operation{OperationID: "update-operation-overrides", Method: "PUT", Path: "/api/operations/{id}/overrides", Summary: "Update operation overrides"},
		func(ctx context.Context, input *UpdateOverridesInput) (*OperationOutput, error) { return nil, nil })

	huma.Register(api, huma.Operation{OperationID: "create-run", Method: "POST", Path: "/api/projects/{project-id}/runs", Summary: "Create and start a run"},
		func(ctx context.Context, input *CreateRunInput) (*RunOutput, error) { return nil, nil })
	huma.Get(api, "/api/projects/{project-id}/runs", func(ctx context.Context, input *ProjectIDParam) (*RunListOutput, error) { return nil, nil })
	huma.Get(api, "/api/runs/{id}", func(ctx context.Context, input *IDParam) (*RunOutput, error) { return nil, nil })
	huma.Get(api, "/api/runs/{id}/result", func(ctx context.Context, input *IDParam) (*RunResultOutput, error) { return nil, nil })
	huma.Get(api, "/api/runs/{id}/report/debug", func(ctx context.Context, input *IDParam) (*DebugReportOutput, error) { return nil, nil })
	huma.Get(api, "/api/runs/{id}/report/ci", func(ctx context.Context, input *IDParam) (*CIReportOutput, error) { return nil, nil })
	huma.Get(api, "/api/runs/{id}/report/mermaid", func(ctx context.Context, input *IDParam) (*MermaidOutput, error) { return nil, nil })
	huma.Get(api, "/api/runs/{id}/artifacts", func(ctx context.Context, input *IDParam) (*ArtifactListOutput, error) { return nil, nil })
	huma.Post(api, "/api/runs/{id}/comments", func(ctx context.Context, input *CreateCommentInput) (*CommentOutput, error) { return nil, nil })
	huma.Get(api, "/api/runs/{id}/comments", func(ctx context.Context, input *IDParam) (*CommentListOutput, error) { return nil, nil })

	b, _ := json.MarshalIndent(api.OpenAPI(), "", "  ")
	if len(os.Args) > 1 {
		os.WriteFile(os.Args[1], b, 0644)
		fmt.Fprintf(os.Stderr, "Written to %s\n", os.Args[1])
	} else {
		fmt.Println(string(b))
	}
}
