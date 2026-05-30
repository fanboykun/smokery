// Package http is the HTTP delivery layer. It registers huma operations
// against an API instance, delegating all business logic to app services.
//
// This is the only package allowed to import HTTP frameworks (echo, huma,
// gorilla/websocket).
package http

import (
	"time"

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

type ProjectWithStatsListOutput struct {
	Body []model.ProjectWithStats
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

type CanvasOperationListOutput struct {
	Body []spec.OperationInfo
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

type RunMetaInput struct {
	IDParam
	Host            string `header:"Host" hidden:"true"`
	ForwardedHost   string `header:"X-Forwarded-Host" hidden:"true"`
	ForwardedProto  string `header:"X-Forwarded-Proto" hidden:"true"`
	ForwardedScheme string `header:"X-Forwarded-Scheme" hidden:"true"`
}

type RunMeta struct {
	ID                     string     `json:"id"`
	ProjectID              string     `json:"project_id"`
	PlanID                 *string    `json:"plan_id,omitempty"`
	Status                 string     `json:"status"`
	CreatedAt              time.Time  `json:"created_at"`
	StartedAt              *time.Time `json:"started_at,omitempty"`
	FinishedAt             *time.Time `json:"finished_at,omitempty"`
	WebSocketURL           string     `json:"websocket_url"`
	FallbackPollIntervalMs int        `json:"fallback_poll_interval_ms"`
	ExpiresAt              time.Time  `json:"expires_at"`
}

type RunMetaOutput struct {
	Body RunMeta
}

type RunListOutput struct {
	Body []model.Run
}

type RunResultOutput struct {
	Body model.StoredRunResult
}

type RunResultDetailsOutput struct {
	Body model.RunResult
}

// --- Plan preview ---

type Diagnostic struct {
	Code       string `json:"code"`
	Severity   string `json:"severity"`
	Message    string `json:"message"`
	EntityID   string `json:"entity_id,omitempty"`
	EntityType string `json:"entity_type,omitempty"`
	Location   string `json:"location,omitempty"`
}

type DiagnosticSummary struct {
	TotalErrors   int  `json:"total_errors"`
	TotalWarnings int  `json:"total_warnings"`
	IsCompilable  bool `json:"is_compilable"`
}

type Diagnostics struct {
	Errors   []Diagnostic      `json:"errors,omitempty"`
	Warnings []Diagnostic      `json:"warnings,omitempty"`
	Summary  DiagnosticSummary `json:"summary"`
}

type PlanPreviewResponse struct {
	Plan        *model.SmokePlan `json:"plan,omitempty"`
	Diagnostics Diagnostics      `json:"diagnostics"`
}

type PlanPreviewOutput struct {
	Body PlanPreviewResponse
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

type ContractReportOutput struct {
	Body report.ContractView
}

type AnalystReportOutput struct {
	Body report.AnalystView
}

type QAReportOutput struct {
	Body report.QAView
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

// --- Failure Classification ---

type ClassifyRunInput struct {
	IDParam
	Body struct {
		Classification string `json:"classification" minLength:"1" doc:"Failure classification"`
		Assignee       string `json:"assignee" doc:"Assigned user"`
		Note           string `json:"note" doc:"Optional note"`
		Author         string `json:"author" minLength:"1" doc:"Who classified"`
	}
}

type FailureClassificationOutput struct {
	Body model.FailureClassification
}

type FailureClassificationListOutput struct {
	Body []model.FailureClassification
}
