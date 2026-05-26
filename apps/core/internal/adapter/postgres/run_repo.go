package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/fanboykun/smokery/apps/core/internal/adapter/postgres/db"
	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

type RunRepo struct{ q *db.Queries }

func NewRunRepo(pool *pgxpool.Pool) *RunRepo {
	return &RunRepo{q: db.New(pool)}
}

var _ port.RunRepo = (*RunRepo)(nil)

func (r *RunRepo) Create(ctx context.Context, projectID uuid.UUID, planID *uuid.UUID) (*model.Run, error) {
	run, err := r.q.CreateRun(ctx, db.CreateRunParams{
		ProjectID: toPgUUID(projectID),
		PlanID:    toPgUUIDPtr(planID),
	})
	if err != nil {
		return nil, err
	}
	return runToModel(run), nil
}

func (r *RunRepo) Get(ctx context.Context, id uuid.UUID) (*model.Run, error) {
	run, err := r.q.GetRun(ctx, toPgUUID(id))
	if err != nil {
		return nil, err
	}
	return runToModel(run), nil
}

func (r *RunRepo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]model.Run, error) {
	runs, err := r.q.ListRunsByProject(ctx, toPgUUID(projectID))
	if err != nil {
		return nil, err
	}
	out := make([]model.Run, len(runs))
	for i, run := range runs {
		out[i] = *runToModel(run)
	}
	return out, nil
}

func (r *RunRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string, startedAt, finishedAt *time.Time) (*model.Run, error) {
	run, err := r.q.UpdateRunStatus(ctx, db.UpdateRunStatusParams{
		ID:         toPgUUID(id),
		Status:     status,
		StartedAt:  toPgTimestamptz(startedAt),
		FinishedAt: toPgTimestamptz(finishedAt),
	})
	if err != nil {
		return nil, err
	}
	return runToModel(run), nil
}

func (r *RunRepo) SaveResult(ctx context.Context, runID uuid.UUID, result []byte) (*model.StoredRunResult, error) {
	rr, err := r.q.CreateRunResult(ctx, db.CreateRunResultParams{
		RunID:  toPgUUID(runID),
		Result: result,
	})
	if err != nil {
		return nil, err
	}
	return runResultToModel(rr), nil
}

func (r *RunRepo) GetResult(ctx context.Context, runID uuid.UUID) (*model.StoredRunResult, error) {
	rr, err := r.q.GetRunResult(ctx, toPgUUID(runID))
	if err != nil {
		return nil, err
	}
	return runResultToModel(rr), nil
}

func runToModel(r db.Run) *model.Run {
	return &model.Run{
		ID:         fromPgUUID(r.ID),
		ProjectID:  fromPgUUID(r.ProjectID),
		PlanID:     fromPgUUIDPtr(r.PlanID),
		Status:     r.Status,
		StartedAt:  fromPgTimestamptzPtr(r.StartedAt),
		FinishedAt: fromPgTimestamptzPtr(r.FinishedAt),
		CreatedAt:  fromPgTimestamptz(r.CreatedAt),
	}
}

func runResultToModel(r db.RunResult) *model.StoredRunResult {
	return &model.StoredRunResult{
		ID:        fromPgUUID(r.ID),
		RunID:     fromPgUUID(r.RunID),
		Result:    r.Result,
		CreatedAt: fromPgTimestamptz(r.CreatedAt),
	}
}
