// Package http is the HTTP delivery layer. It registers huma operations
// against an API instance, delegating all business logic to app services.
//
// This is the only package allowed to import HTTP frameworks (echo, huma,
// gorilla/websocket).
package http

import (
	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/report"
	"github.com/fanboykun/smokery/apps/core/internal/spec"
)

// --- Path params ---

type IDParam struct {
	ID string `path:"id" format:"uuid" doc:"Resource UUID"`
}

type ProjectIDParam struct {
	ProjectID string `path:"project-id" format:"uuid" doc:"Project UUID"`
}

type SpecIDParam struct {
	SpecID string `path:"spec-id" format:"uuid" doc:"Spec UUID"`
}

// --- Health ---

type HealthOutput struct {
	Body struct {
		Status string `json:"status" example:"ok" doc:"Service health status"`
	}
}

// --- Projects ---

type CreateProjectInput struct {
	Body struct {
		Name        string `json:"name" minLength:"1" doc:"Project name"`
		Description string `json:"description" doc:"Project description"`
	}
}

type UpdateProjectInput struct {
	IDParam
	Body struct {
		Name        string `json:"name" minLength:"1" doc:"Project name"`
		Description string `json:"description" doc:"Project description"`
	}
}

type ProjectOutput struct {
	Body model.Project
}

type ProjectListOutput struct {
	Body []model.Project
}

// --- Specs ---

type ImportSpecInput struct {
	ProjectIDParam
	RawBody []byte
}

type SpecAnalysisOutput struct {
	Body spec.Analysis
}

// --- Operations ---

type OperationListOutput struct {
	Body []model.Operation
}

type OperationOutput struct {
	Body model.Operation
}

type UpdateClassificationInput struct {
	IDParam
	Body struct {
		Classification string `json:"classification" minLength:"1" doc:"Operation classification"`
		IsDestructive  bool   `json:"is_destructive" doc:"Whether operation is destructive"`
	}
}

type UpdateOverridesInput struct {
	IDParam
	RawBody []byte
}

// --- Runs ---

type CreateRunInput struct {
	ProjectIDParam
	Body struct {
		PlanID string           `json:"plan_id" doc:"Plan UUID"`
		Plan   *model.SmokePlan `json:"plan" doc:"Compiled smoke plan to execute"`
	}
}

type RunOutput struct {
	Body model.Run
}

type RunListOutput struct {
	Body []model.Run
}

type RunResultOutput struct {
	Body model.StoredRunResult
}

// --- Reports ---

type DebugReportOutput struct {
	Body report.DebugView
}

type CIReportOutput struct {
	Body report.CISummary
}

type MermaidOutput struct {
	Body string
}

// --- Artifacts ---

type ArtifactListOutput struct {
	Body []model.Artifact
}

// --- Comments ---

type CreateCommentInput struct {
	IDParam
	Body struct {
		Author string `json:"author" minLength:"1" doc:"Comment author"`
		Body   string `json:"body" minLength:"1" doc:"Comment body"`
	}
}

type CommentOutput struct {
	Body model.Comment
}

type CommentListOutput struct {
	Body []model.Comment
}
