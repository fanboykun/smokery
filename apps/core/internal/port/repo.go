package port

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
)

type ProjectRepo interface {
	Create(ctx context.Context, name, description string) (*model.Project, error)
	Get(ctx context.Context, id uuid.UUID) (*model.Project, error)
	List(ctx context.Context) ([]model.Project, error)
	Update(ctx context.Context, id uuid.UUID, name, description string) (*model.Project, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type SpecRepo interface {
	Create(ctx context.Context, in model.Spec) (*model.Spec, error)
	Get(ctx context.Context, id uuid.UUID) (*model.Spec, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]model.Spec, error)
}

type OperationRepo interface {
	Create(ctx context.Context, in model.Operation) (*model.Operation, error)
	ListBySpec(ctx context.Context, specID uuid.UUID) ([]model.Operation, error)
	UpdateClassification(ctx context.Context, id uuid.UUID, classification string, isDestructive bool) (*model.Operation, error)
	UpdateOverrides(ctx context.Context, id uuid.UUID, overrides []byte) (*model.Operation, error)
}

type RunRepo interface {
	Create(ctx context.Context, projectID uuid.UUID, planID *uuid.UUID) (*model.Run, error)
	Get(ctx context.Context, id uuid.UUID) (*model.Run, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]model.Run, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, startedAt, finishedAt *time.Time) (*model.Run, error)
	SaveResult(ctx context.Context, runID uuid.UUID, result []byte) (*model.StoredRunResult, error)
	GetResult(ctx context.Context, runID uuid.UUID) (*model.StoredRunResult, error)
}

type CommentRepo interface {
	Create(ctx context.Context, runID uuid.UUID, author, body string) (*model.Comment, error)
	ListByRun(ctx context.Context, runID uuid.UUID) ([]model.Comment, error)
}

type ArtifactRepo interface {
	Create(ctx context.Context, runID uuid.UUID, typ, path string) (*model.Artifact, error)
	ListByRun(ctx context.Context, runID uuid.UUID) ([]model.Artifact, error)
}
