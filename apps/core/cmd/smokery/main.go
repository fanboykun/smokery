// Command smokery is the standalone CLI smoke runner. It uses the same
// compiler, runner, and app services as the API server, but with in-memory
// repos and a filesystem blob store. No database required.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/adapter/memory"
	"github.com/fanboykun/smokery/apps/core/internal/app"
	"github.com/fanboykun/smokery/apps/core/internal/delivery/cli"
	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
	"github.com/fanboykun/smokery/apps/core/internal/runner"
)

// noopJobs is a no-op JobEnqueuer for the CLI. The CLI runs plans
// synchronously via the runner; it does not need async job orchestration.
type noopJobs struct{}

func (noopJobs) EnqueueRun(context.Context, uuid.UUID, *model.SmokePlan) error { return nil }

var _ port.JobEnqueuer = noopJobs{}

func main() {
	// Adapters
	store := memory.NewStore()
	projectRepo := memory.NewProjectRepo(store)
	specRepo := memory.NewSpecRepo(store)
	operationRepo := memory.NewOperationRepo(store)
	runRepo := memory.NewRunRepo(store)

	// Runner (no event emitter for CLI)
	rnr := runner.New(runner.DefaultOptions())

	// App services
	svcs := &cli.Services{
		Project:   app.NewProjectService(projectRepo),
		Spec:      app.NewSpecService(specRepo, operationRepo),
		Operation: app.NewOperationService(operationRepo),
		Run:       app.NewRunService(runRepo, &noopJobs{}),
		Report:    app.NewReportService(runRepo),
		Runner:    rnr,
	}

	if err := cli.NewRootCmd(svcs).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
