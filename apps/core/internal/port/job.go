package port

import (
	"context"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
)

// JobEnqueuer enqueues a run for asynchronous execution.
type JobEnqueuer interface {
	EnqueueRun(ctx context.Context, runID uuid.UUID, plan *model.SmokePlan) error
	CancelRun(ctx context.Context, runID uuid.UUID) error
}
