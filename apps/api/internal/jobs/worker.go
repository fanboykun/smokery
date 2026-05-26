package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/fanboykun/smokery/apps/api/internal/db"
	"github.com/fanboykun/smokery/apps/api/internal/model"
	"github.com/fanboykun/smokery/apps/api/internal/runner"
)

type EventType string

const (
	EventRunStarted  EventType = "run_started"
	EventStepDone    EventType = "step_done"
	EventRunFinished EventType = "run_finished"
)

type RunEvent struct {
	Type  EventType `json:"type"`
	RunID string    `json:"run_id"`
	Data  any       `json:"data,omitempty"`
}

type Worker struct {
	pool    *pgxpool.Pool
	queries *db.Queries
	runner  *runner.Runner
	mu      sync.RWMutex
	subs    map[string][]chan RunEvent
}

func NewWorker(pool *pgxpool.Pool) *Worker {
	return &Worker{
		pool:    pool,
		queries: db.New(pool),
		runner:  runner.New(),
		subs:    make(map[string][]chan RunEvent),
	}
}

func (w *Worker) Subscribe(runID string) chan RunEvent {
	ch := make(chan RunEvent, 64)
	w.mu.Lock()
	w.subs[runID] = append(w.subs[runID], ch)
	w.mu.Unlock()
	return ch
}

func (w *Worker) Unsubscribe(runID string, ch chan RunEvent) {
	w.mu.Lock()
	defer w.mu.Unlock()
	subs := w.subs[runID]
	for i, s := range subs {
		if s == ch {
			w.subs[runID] = append(subs[:i], subs[i+1:]...)
			close(ch)
			return
		}
	}
}

func (w *Worker) broadcast(runID string, event RunEvent) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	for _, ch := range w.subs[runID] {
		select {
		case ch <- event:
		default:
		}
	}
}

func (w *Worker) Enqueue(ctx context.Context, runID pgtype.UUID, plan *model.SmokePlan) {
	go w.execute(ctx, runID, plan)
}

func (w *Worker) execute(ctx context.Context, runID pgtype.UUID, plan *model.SmokePlan) {
	now := pgtype.Timestamptz{Time: time.Now(), Valid: true}
	w.queries.UpdateRunStatus(ctx, db.UpdateRunStatusParams{
		ID: runID, Status: "running", StartedAt: now, FinishedAt: pgtype.Timestamptz{},
	})

	runIDStr := uuidToString(runID)
	w.broadcast(runIDStr, RunEvent{Type: EventRunStarted, RunID: runIDStr})

	result := w.runner.Execute(ctx, plan)

	finished := pgtype.Timestamptz{Time: time.Now(), Valid: true}
	status := "completed"
	if result.Status == "failed" {
		status = "failed"
	}
	w.queries.UpdateRunStatus(ctx, db.UpdateRunStatusParams{
		ID: runID, Status: status, StartedAt: now, FinishedAt: finished,
	})

	resultJSON, _ := json.Marshal(result)
	w.queries.CreateRunResult(ctx, db.CreateRunResultParams{RunID: runID, Result: resultJSON})

	w.broadcast(runIDStr, RunEvent{Type: EventRunFinished, RunID: runIDStr, Data: result})
}

func uuidToString(id pgtype.UUID) string {
	b := id.Bytes
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
