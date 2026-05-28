package inproc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
	"github.com/fanboykun/smokery/apps/core/internal/report"
	"github.com/fanboykun/smokery/apps/core/internal/runner"
)

// Worker is an in-process job runner implementing port.JobEnqueuer.
type Worker struct {
	runRepo   port.RunRepo
	bus       port.EventBus
	runner    *runner.Runner
	blob      port.BlobStore
	artifacts port.ArtifactRepo

	mu      sync.Mutex
	cancels map[uuid.UUID]context.CancelFunc
}

func NewWorker(runRepo port.RunRepo, bus port.EventBus, r *runner.Runner) *Worker {
	return &Worker{runRepo: runRepo, bus: bus, runner: r, cancels: make(map[uuid.UUID]context.CancelFunc)}
}

// WithArtifacts enables artifact persistence after runs complete.
func (w *Worker) WithArtifacts(blob port.BlobStore, artifacts port.ArtifactRepo) *Worker {
	w.blob = blob
	w.artifacts = artifacts
	return w
}

var _ port.JobEnqueuer = (*Worker)(nil)

func (w *Worker) EnqueueRun(ctx context.Context, runID uuid.UUID, plan *model.SmokePlan) error {
	runCtx, cancel := context.WithCancel(context.Background())
	w.mu.Lock()
	w.cancels[runID] = cancel
	w.mu.Unlock()
	go w.execute(runCtx, runID, plan)
	return nil
}

func (w *Worker) CancelRun(_ context.Context, runID uuid.UUID) error {
	w.mu.Lock()
	cancel, ok := w.cancels[runID]
	w.mu.Unlock()
	if ok {
		cancel()
	}
	return nil
}

func (w *Worker) execute(ctx context.Context, runID uuid.UUID, plan *model.SmokePlan) {
	defer func() {
		w.mu.Lock()
		delete(w.cancels, runID)
		w.mu.Unlock()
	}()

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
	if ctx.Err() != nil {
		status = "cancelled"
		result.Status = "cancelled"
	} else if result.Status == "failed" || result.Status == "error" {
		status = "failed"
	}
	_, _ = w.runRepo.UpdateStatus(ctx, runID, status, &now, &finished)

	resultJSON, _ := json.Marshal(result)
	_, _ = w.runRepo.SaveResult(ctx, runID, resultJSON)

	// Persist artifacts if blob store is configured
	if w.blob != nil && w.artifacts != nil {
		w.persistArtifacts(context.Background(), runID, result)
	}

	w.bus.Publish(port.Event{Type: port.EventRunFinished, RunID: runID.String(), Data: result})
}

func (w *Worker) persistArtifacts(ctx context.Context, runID uuid.UUID, result *model.RunResult) {
	rid := runID.String()

	// JSON artifact
	jsonKey := fmt.Sprintf("runs/%s/result.json", rid)
	jsonData := report.GenerateJSONArtifact(result)
	if err := w.blob.Put(ctx, jsonKey, bytes.NewReader(jsonData), "application/json"); err == nil {
		_, _ = w.artifacts.Create(ctx, runID, "json_report", jsonKey)
	}

	// HTML artifact
	htmlKey := fmt.Sprintf("runs/%s/report.html", rid)
	htmlData, err := report.GenerateHTMLReport(result)
	if err == nil {
		if err := w.blob.Put(ctx, htmlKey, bytes.NewReader(htmlData), "text/html"); err == nil {
			_, _ = w.artifacts.Create(ctx, runID, "html_report", htmlKey)
		}
	}
}
