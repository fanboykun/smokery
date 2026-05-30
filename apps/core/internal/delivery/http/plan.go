package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/app"
	"github.com/fanboykun/smokery/apps/core/internal/model"
)

type PlanPreviewInput struct {
	ProjectIDParam
	Body model.ProjectConfig
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

		diagnostics := Diagnostics{
			Errors:   make([]Diagnostic, 0, len(out.Errors)),
			Warnings: make([]Diagnostic, 0, len(out.Warnings)),
		}
		for _, err := range out.Errors {
			diagnostics.Errors = append(diagnostics.Errors, diagnosticFromCompileError(err))
		}
		for _, warning := range out.Warnings {
			diagnostics.Warnings = append(diagnostics.Warnings, diagnosticFromCompileError(warning))
		}
		diagnostics.Summary = DiagnosticSummary{
			TotalErrors:   len(diagnostics.Errors),
			TotalWarnings: len(diagnostics.Warnings),
			IsCompilable:  out.Plan != nil && len(diagnostics.Errors) == 0,
		}

		return &PlanPreviewOutput{Body: PlanPreviewResponse{
			Plan:        out.Plan,
			Diagnostics: diagnostics,
		}}, nil
	})
}

func diagnosticFromCompileError(err model.CompileError) Diagnostic {
	return Diagnostic{
		Code:       diagnosticCode(err),
		Severity:   err.Severity,
		Message:    err.Message,
		EntityID:   err.Entity,
		EntityType: err.Stage,
		Location:   err.Path,
	}
}

func diagnosticCode(err model.CompileError) string {
	message := strings.ToLower(err.Message)
	path := strings.ToLower(err.Path)
	stage := strings.ToLower(err.Stage)

	switch {
	case strings.Contains(message, "auth"):
		return "MISSING_AUTH_PROFILE"
	case strings.Contains(message, "operation") && strings.Contains(message, "not found"):
		return "MISSING_OPERATION"
	case strings.Contains(stage, "suite") && strings.Contains(path, "selector") && strings.Contains(message, "no operations"):
		return "INVALID_SELECTOR"
	case strings.Contains(message, "unused") || strings.Contains(message, "not used"):
		return "UNUSED_OPERATION"
	case strings.Contains(message, "assertion"):
		return "INVALID_ASSERTION"
	case strings.Contains(message, "schema"):
		return "SCHEMA_MISMATCH"
	default:
		return "COMPILATION_ERROR"
	}
}
