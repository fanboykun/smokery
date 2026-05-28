// Command openapi generates the OpenAPI 3.1 spec without needing a database.
// It registers the same huma operations as the server but with stub services
// that satisfy the app.* signatures using nil ports. Only the schemas/paths
// matter for spec generation, not runtime behavior.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/fanboykun/smokery/apps/core/internal/app"
	deliveryhttp "github.com/fanboykun/smokery/apps/core/internal/delivery/http"
	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

func main() {
	e := echo.New()
	cfg := huma.DefaultConfig("Smokery API", "1.0.0")
	cfg.Servers = []*huma.Server{{URL: "http://localhost:8080"}}
	api := humaecho.New(e, cfg)

	// Register all the same operations as the real server. We pass services
	// backed by no-op port implementations so signatures resolve without a DB.
	projectSvc := app.NewProjectService(noopProjectRepo{})
	specSvc := app.NewSpecService(noopSpecRepo{}, noopOperationRepo{})
	operationSvc := app.NewOperationService(noopOperationRepo{})
	runSvc := app.NewRunService(noopRunRepo{}, noopJobs{})
	reportSvc := app.NewReportService(noopRunRepo{})
	commentSvc := app.NewCommentService(noopCommentRepo{})
	artifactSvc := app.NewArtifactService(noopArtifactRepo{})

	deliveryhttp.RegisterHealthCheck(api, noopHealthChecker{})
	deliveryhttp.RegisterProjects(api, projectSvc)
	deliveryhttp.RegisterSpecs(api, specSvc)
	deliveryhttp.RegisterOperations(api, operationSvc)
	deliveryhttp.RegisterRuns(api, runSvc)
	deliveryhttp.RegisterReports(api, reportSvc)
	deliveryhttp.RegisterComments(api, commentSvc)
	deliveryhttp.RegisterArtifacts(api, artifactSvc)
	deliveryhttp.RegisterPlan(api, app.NewPlanService(noopSpecRepo{}, noopOperationRepo{}))

	b, _ := json.MarshalIndent(api.OpenAPI(), "", "  ")
	if len(os.Args) > 1 {
		_ = os.WriteFile(os.Args[1], b, 0644)
		fmt.Fprintf(os.Stderr, "Written to %s\n", os.Args[1])
	} else {
		fmt.Println(string(b))
	}
}

// --- No-op port implementations (only used by the openapi tool) ---

type noopProjectRepo struct{}

func (noopProjectRepo) Create(context.Context, string, string) (*model.Project, error) {
	return nil, nil
}
func (noopProjectRepo) Get(context.Context, uuid.UUID) (*model.Project, error) { return nil, nil }
func (noopProjectRepo) List(context.Context) ([]model.Project, error)          { return nil, nil }
func (noopProjectRepo) Update(context.Context, uuid.UUID, string, string) (*model.Project, error) {
	return nil, nil
}
func (noopProjectRepo) Delete(context.Context, uuid.UUID) error { return nil }

type noopSpecRepo struct{}

func (noopSpecRepo) Create(context.Context, model.Spec) (*model.Spec, error)         { return nil, nil }
func (noopSpecRepo) Get(context.Context, uuid.UUID) (*model.Spec, error)             { return nil, nil }
func (noopSpecRepo) ListByProject(context.Context, uuid.UUID) ([]model.Spec, error)  { return nil, nil }

type noopOperationRepo struct{}

func (noopOperationRepo) Create(context.Context, model.Operation) (*model.Operation, error) {
	return nil, nil
}
func (noopOperationRepo) ListBySpec(context.Context, uuid.UUID) ([]model.Operation, error) {
	return nil, nil
}
func (noopOperationRepo) UpdateClassification(context.Context, uuid.UUID, string, bool) (*model.Operation, error) {
	return nil, nil
}
func (noopOperationRepo) UpdateOverrides(context.Context, uuid.UUID, []byte) (*model.Operation, error) {
	return nil, nil
}

type noopRunRepo struct{}

func (noopRunRepo) Create(context.Context, uuid.UUID, *uuid.UUID) (*model.Run, error) {
	return nil, nil
}
func (noopRunRepo) Get(context.Context, uuid.UUID) (*model.Run, error)             { return nil, nil }
func (noopRunRepo) ListByProject(context.Context, uuid.UUID) ([]model.Run, error)  { return nil, nil }
func (noopRunRepo) UpdateStatus(context.Context, uuid.UUID, string, *time.Time, *time.Time) (*model.Run, error) {
	return nil, nil
}
func (noopRunRepo) SaveResult(context.Context, uuid.UUID, []byte) (*model.StoredRunResult, error) {
	return nil, nil
}
func (noopRunRepo) GetResult(context.Context, uuid.UUID) (*model.StoredRunResult, error) {
	return nil, nil
}
func (noopRunRepo) DeleteOlderThan(context.Context, time.Time) (int, error) { return 0, nil }

type noopCommentRepo struct{}

func (noopCommentRepo) Create(context.Context, uuid.UUID, string, string) (*model.Comment, error) {
	return nil, nil
}
func (noopCommentRepo) ListByRun(context.Context, uuid.UUID) ([]model.Comment, error) {
	return nil, nil
}

type noopArtifactRepo struct{}

func (noopArtifactRepo) Create(context.Context, uuid.UUID, string, string) (*model.Artifact, error) {
	return nil, nil
}
func (noopArtifactRepo) ListByRun(context.Context, uuid.UUID) ([]model.Artifact, error) {
	return nil, nil
}
func (noopArtifactRepo) DeleteByRun(context.Context, uuid.UUID) ([]model.Artifact, error) {
	return nil, nil
}

type noopJobs struct{}

func (noopJobs) EnqueueRun(context.Context, uuid.UUID, *model.SmokePlan) error { return nil }
func (noopJobs) CancelRun(context.Context, uuid.UUID) error                    { return nil }

type noopHealthChecker struct{}

func (noopHealthChecker) PingDB(context.Context) error   { return nil }
func (noopHealthChecker) PingBlob(context.Context) error { return nil }

// Compile-time interface checks
var (
	_ port.ProjectRepo   = noopProjectRepo{}
	_ port.SpecRepo      = noopSpecRepo{}
	_ port.OperationRepo = noopOperationRepo{}
	_ port.RunRepo       = noopRunRepo{}
	_ port.CommentRepo   = noopCommentRepo{}
	_ port.ArtifactRepo  = noopArtifactRepo{}
	_ port.JobEnqueuer   = noopJobs{}
	_ deliveryhttp.HealthChecker = noopHealthChecker{}
)
