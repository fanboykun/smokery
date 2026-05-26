package inproc

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
	"github.com/fanboykun/smokery/apps/core/internal/runner"
)

// Worker is an in-process job runner implementing port.JobEnqueuer.
type Worker struct {
	runRepo port.RunRepo
	bus     port.EventBus
	runner  *runner.Runner
}

func NewWorker(runRepo port.RunRepo, bus port.EventBus, r *runner.Runner) *Worker {
	return &Worker{runRepo: runRepo, bus: bus, runner: r}
}

var _ port.JobEnqueuer = (*Worker)(nil)

func (w *Worker) EnqueueRun(ctx context.Context, runID uuid.UUID, plan *model.SmokePlan) error {
	go w.execute(context.Background(), runID, plan)
	return nil
}

func (w *Worker) execute(ctx context.Context, runID uuid.UUID, plan *model.SmokePlan) {
	now := time.Now()
	_, _ = w.runRepo.UpdateStatus(ctx, runID, "running", &now, nil)
	w.bus.Publish(port.Event{Type: port.EventRunStarted, RunID: runID.String()})

	var result *model.RunResult
	if plan != nil {
		result = w.runner.Execute(ctx, plan)
	} else {
		result = &model.RunResult{RunID: runID.String(), Status: "failed", StartedAt: now, FinishedAt: time.Now()}
	}

	finished := time.Now()
	status := "completed"
	if result.Status == "failed" || result.Status == "error" {
		status = "failed"
	}
	_, _ = w.runRepo.UpdateStatus(ctx, runID, status, &now, &finished)

	resultJSON, _ := json.Marshal(result)
	_, _ = w.runRepo.SaveResult(ctx, runID, resultJSON)

	w.bus.Publish(port.Event{Type: port.EventRunFinished, RunID: runID.String(), Data: result})
}
