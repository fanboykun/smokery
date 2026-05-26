package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/fanboykun/smokery/apps/api/internal/config"
	"github.com/fanboykun/smokery/apps/api/internal/db"
	"github.com/fanboykun/smokery/apps/api/internal/jobs"
	"github.com/fanboykun/smokery/apps/api/internal/model"
	"github.com/fanboykun/smokery/apps/api/internal/report"
	"github.com/fanboykun/smokery/apps/api/internal/spec"
)

// --- Input/Output types for huma operations ---

type CreateProjectInput struct {
	Body struct {
		Name        string `json:"name" minLength:"1" doc:"Project name"`
		Description string `json:"description" doc:"Project description"`
	}
}

type ProjectOutput struct {
	Body db.Project
}

type ProjectListOutput struct {
	Body []db.Project
}

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

type SpecAnalysisOutput struct {
	Body spec.Analysis
}

type OperationListOutput struct {
	Body []db.Operation
}

type UpdateClassificationInput struct {
	IDParam
	Body struct {
		Classification string `json:"classification" minLength:"1" doc:"Operation classification"`
		IsDestructive  bool   `json:"is_destructive" doc:"Whether operation is destructive"`
	}
}

type OperationOutput struct {
	Body db.Operation
}

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

type RunOutput struct {
	Body db.Run
}

type RunListOutput struct {
	Body []db.Run
}

type RunResultOutput struct {
	Body db.RunResult
}

type DebugReportOutput struct {
	Body report.DebugView
}

type CIReportOutput struct {
	Body report.CISummary
}

type MermaidOutput struct {
	Body string
}

type ArtifactListOutput struct {
	Body []db.Artifact
}

type CreateCommentInput struct {
	IDParam
	Body struct {
		Author string `json:"author" minLength:"1" doc:"Comment author"`
		Body   string `json:"body" minLength:"1" doc:"Comment body"`
	}
}

type CommentOutput struct {
	Body db.Comment
}

type CommentListOutput struct {
	Body []db.Comment
}

type HealthOutput struct {
	Body struct {
		Status string `json:"status" example:"ok" doc:"Service health status"`
	}
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer pool.Close()

	queries := db.New(pool)
	worker := jobs.NewWorker(pool)

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
	}))

	// Create huma API with OpenAPI config
	humaConfig := huma.DefaultConfig("Smokery API", "1.0.0")
	humaConfig.Servers = []*huma.Server{{URL: "http://localhost:" + cfg.Port}}
	api := humaecho.New(e, humaConfig)

	// Health
	huma.Get(api, "/health", func(ctx context.Context, input *struct{}) (*HealthOutput, error) {
		return &HealthOutput{Body: struct {
			Status string `json:"status" example:"ok" doc:"Service health status"`
		}{Status: "ok"}}, nil
	})

	// Projects
	huma.Post(api, "/api/projects", func(ctx context.Context, input *CreateProjectInput) (*ProjectOutput, error) {
		p, err := queries.CreateProject(ctx, db.CreateProjectParams{Name: input.Body.Name, Description: input.Body.Description})
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to create project", err)
		}
		return &ProjectOutput{Body: p}, nil
	})

	huma.Get(api, "/api/projects", func(ctx context.Context, input *struct{}) (*ProjectListOutput, error) {
		projects, err := queries.ListProjects(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list projects", err)
		}
		return &ProjectListOutput{Body: projects}, nil
	})

	huma.Get(api, "/api/projects/{id}", func(ctx context.Context, input *IDParam) (*ProjectOutput, error) {
		id, err := parseUUID(input.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id")
		}
		p, err := queries.GetProject(ctx, id)
		if err != nil {
			return nil, huma.Error404NotFound("project not found")
		}
		return &ProjectOutput{Body: p}, nil
	})

	huma.Put(api, "/api/projects/{id}", func(ctx context.Context, input *UpdateProjectInput) (*ProjectOutput, error) {
		id, err := parseUUID(input.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id")
		}
		p, err := queries.UpdateProject(ctx, db.UpdateProjectParams{ID: id, Name: input.Body.Name, Description: input.Body.Description})
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to update project", err)
		}
		return &ProjectOutput{Body: p}, nil
	})

	huma.Delete(api, "/api/projects/{id}", func(ctx context.Context, input *IDParam) (*struct{}, error) {
		id, err := parseUUID(input.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id")
		}
		if err := queries.DeleteProject(ctx, id); err != nil {
			return nil, huma.Error500InternalServerError("failed to delete project", err)
		}
		return nil, nil
	})

	// Spec import
	huma.Register(api, huma.Operation{
		OperationID: "import-spec",
		Method:      http.MethodPost,
		Path:        "/api/projects/{project-id}/specs",
		Summary:     "Import OpenAPI spec",
	}, func(ctx context.Context, input *ImportSpecInput) (*SpecAnalysisOutput, error) {
		projectID, err := parseUUID(input.ProjectID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid project id")
		}
		analysis, err := spec.Parse(input.RawBody)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity("failed to parse spec", err)
		}
		analysisJSON, _ := json.Marshal(analysis)
		s, err := queries.CreateSpec(ctx, db.CreateSpecParams{
			ProjectID: projectID, Version: analysis.Version, Title: analysis.Title,
			Raw: input.RawBody, Analysis: analysisJSON,
		})
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to store spec", err)
		}
		for _, op := range analysis.Operations {
			tagsJSON, _ := json.Marshal(op.Tags)
			queries.CreateOperation(ctx, db.CreateOperationParams{
				SpecID: s.ID, OperationID: op.OperationID, Method: op.Method,
				Path: op.Path, Summary: op.Summary, Tags: tagsJSON,
				Classification: op.Classification, IsDestructive: op.IsDestructive,
			})
		}
		return &SpecAnalysisOutput{Body: *analysis}, nil
	})

	// Operations
	huma.Get(api, "/api/specs/{spec-id}/operations", func(ctx context.Context, input *SpecIDParam) (*OperationListOutput, error) {
		specID, err := parseUUID(input.SpecID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid spec id")
		}
		ops, err := queries.ListOperationsBySpec(ctx, specID)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list operations", err)
		}
		return &OperationListOutput{Body: ops}, nil
	})

	huma.Put(api, "/api/operations/{id}/classification", func(ctx context.Context, input *UpdateClassificationInput) (*OperationOutput, error) {
		id, err := parseUUID(input.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id")
		}
		op, err := queries.UpdateOperationClassification(ctx, db.UpdateOperationClassificationParams{
			ID: id, Classification: input.Body.Classification, IsDestructive: input.Body.IsDestructive,
		})
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to update classification", err)
		}
		return &OperationOutput{Body: op}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-operation-overrides",
		Method:      http.MethodPut,
		Path:        "/api/operations/{id}/overrides",
		Summary:     "Update operation overrides",
	}, func(ctx context.Context, input *UpdateOverridesInput) (*OperationOutput, error) {
		id, err := parseUUID(input.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id")
		}
		op, err := queries.UpdateOperationOverrides(ctx, db.UpdateOperationOverridesParams{ID: id, Overrides: input.RawBody})
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to update overrides", err)
		}
		return &OperationOutput{Body: op}, nil
	})

	// Runs
	huma.Register(api, huma.Operation{
		OperationID: "create-run",
		Method:      http.MethodPost,
		Path:        "/api/projects/{project-id}/runs",
		Summary:     "Create and start a run",
	}, func(ctx context.Context, input *CreateRunInput) (*RunOutput, error) {
		projectID, err := parseUUID(input.ProjectID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid project id")
		}
		planID, _ := parseUUID(input.Body.PlanID)
		run, err := queries.CreateRun(ctx, db.CreateRunParams{ProjectID: projectID, PlanID: planID})
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to create run", err)
		}
		worker.Enqueue(context.Background(), run.ID, input.Body.Plan)
		return &RunOutput{Body: run}, nil
	})

	huma.Get(api, "/api/projects/{project-id}/runs", func(ctx context.Context, input *ProjectIDParam) (*RunListOutput, error) {
		projectID, err := parseUUID(input.ProjectID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid project id")
		}
		runs, err := queries.ListRunsByProject(ctx, projectID)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list runs", err)
		}
		return &RunListOutput{Body: runs}, nil
	})

	huma.Get(api, "/api/runs/{id}", func(ctx context.Context, input *IDParam) (*RunOutput, error) {
		id, err := parseUUID(input.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id")
		}
		run, err := queries.GetRun(ctx, id)
		if err != nil {
			return nil, huma.Error404NotFound("run not found")
		}
		return &RunOutput{Body: run}, nil
	})

	huma.Get(api, "/api/runs/{id}/result", func(ctx context.Context, input *IDParam) (*RunResultOutput, error) {
		id, err := parseUUID(input.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id")
		}
		result, err := queries.GetRunResult(ctx, id)
		if err != nil {
			return nil, huma.Error404NotFound("result not found")
		}
		return &RunResultOutput{Body: result}, nil
	})

	// Reports
	huma.Get(api, "/api/runs/{id}/report/debug", func(ctx context.Context, input *IDParam) (*DebugReportOutput, error) {
		id, err := parseUUID(input.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id")
		}
		rr, err := queries.GetRunResult(ctx, id)
		if err != nil {
			return nil, huma.Error404NotFound("result not found")
		}
		var run model.RunResult
		if err := json.Unmarshal(rr.Result, &run); err != nil {
			return nil, huma.Error500InternalServerError("invalid result data", err)
		}
		return &DebugReportOutput{Body: *report.GenerateDebugView(&run)}, nil
	})

	huma.Get(api, "/api/runs/{id}/report/ci", func(ctx context.Context, input *IDParam) (*CIReportOutput, error) {
		id, err := parseUUID(input.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id")
		}
		rr, err := queries.GetRunResult(ctx, id)
		if err != nil {
			return nil, huma.Error404NotFound("result not found")
		}
		var run model.RunResult
		if err := json.Unmarshal(rr.Result, &run); err != nil {
			return nil, huma.Error500InternalServerError("invalid result data", err)
		}
		return &CIReportOutput{Body: *report.GenerateCISummary(&run)}, nil
	})

	huma.Get(api, "/api/runs/{id}/report/mermaid", func(ctx context.Context, input *IDParam) (*MermaidOutput, error) {
		id, err := parseUUID(input.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id")
		}
		rr, err := queries.GetRunResult(ctx, id)
		if err != nil {
			return nil, huma.Error404NotFound("result not found")
		}
		var run model.RunResult
		if err := json.Unmarshal(rr.Result, &run); err != nil {
			return nil, huma.Error500InternalServerError("invalid result data", err)
		}
		return &MermaidOutput{Body: report.GenerateMermaidDiagram(&run)}, nil
	})

	// Artifacts
	huma.Get(api, "/api/runs/{id}/artifacts", func(ctx context.Context, input *IDParam) (*ArtifactListOutput, error) {
		id, err := parseUUID(input.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id")
		}
		artifacts, err := queries.ListArtifactsByRun(ctx, id)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list artifacts", err)
		}
		return &ArtifactListOutput{Body: artifacts}, nil
	})

	// Comments
	huma.Post(api, "/api/runs/{id}/comments", func(ctx context.Context, input *CreateCommentInput) (*CommentOutput, error) {
		id, err := parseUUID(input.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id")
		}
		comment, err := queries.CreateComment(ctx, db.CreateCommentParams{RunID: id, Author: input.Body.Author, Body: input.Body.Body})
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to create comment", err)
		}
		return &CommentOutput{Body: comment}, nil
	})

	huma.Get(api, "/api/runs/{id}/comments", func(ctx context.Context, input *IDParam) (*CommentListOutput, error) {
		id, err := parseUUID(input.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id")
		}
		comments, err := queries.ListCommentsByRun(ctx, id)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list comments", err)
		}
		return &CommentListOutput{Body: comments}, nil
	})

	// WebSocket (stays as raw Echo handler - not part of OpenAPI)
	var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	e.GET("/ws/runs/:id", func(c echo.Context) error {
		runID := c.Param("id")
		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()
		ch := worker.Subscribe(runID)
		defer worker.Unsubscribe(runID, ch)
		for event := range ch {
			if err := ws.WriteJSON(event); err != nil {
				break
			}
			if event.Type == jobs.EventRunFinished {
				break
			}
		}
		return nil
	})

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server error")
		}
	}()
	log.Info().Str("port", cfg.Port).Msg("server started")
	log.Info().Msg("OpenAPI spec at /openapi.json")
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	e.Shutdown(shutdownCtx)
}

func parseUUID(s string) (pgtype.UUID, error) {
	var id pgtype.UUID
	err := id.Scan(s)
	return id, err
}
