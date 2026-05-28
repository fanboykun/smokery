package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/fanboykun/smokery/apps/core/internal/adapter/postgres/db"
	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

type ArtifactRepo struct {
	q    *db.Queries
	pool *pgxpool.Pool
}

func NewArtifactRepo(pool *pgxpool.Pool) *ArtifactRepo {
	return &ArtifactRepo{q: db.New(pool), pool: pool}
}

var _ port.ArtifactRepo = (*ArtifactRepo)(nil)

func (r *ArtifactRepo) Create(ctx context.Context, runID uuid.UUID, typ, path string) (*model.Artifact, error) {
	a, err := r.q.CreateArtifact(ctx, db.CreateArtifactParams{
		RunID: toPgUUID(runID),
		Type:  typ,
		Path:  path,
	})
	if err != nil {
		return nil, err
	}
	return artifactToModel(a), nil
}

func (r *ArtifactRepo) ListByRun(ctx context.Context, runID uuid.UUID) ([]model.Artifact, error) {
	as, err := r.q.ListArtifactsByRun(ctx, toPgUUID(runID))
	if err != nil {
		return nil, err
	}
	out := make([]model.Artifact, len(as))
	for i, a := range as {
		out[i] = *artifactToModel(a)
	}
	return out, nil
}

func artifactToModel(a db.Artifact) *model.Artifact {
	return &model.Artifact{
		ID:        fromPgUUID(a.ID),
		RunID:     fromPgUUID(a.RunID),
		Type:      a.Type,
		Path:      a.Path,
		CreatedAt: fromPgTimestamptz(a.CreatedAt),
	}
}

func (r *ArtifactRepo) DeleteByRun(ctx context.Context, runID uuid.UUID) ([]model.Artifact, error) {
	arts, err := r.ListByRun(ctx, runID)
	if err != nil {
		return nil, err
	}
	_, err = r.pool.Exec(ctx, "DELETE FROM artifacts WHERE run_id = $1", toPgUUID(runID))
	if err != nil {
		return nil, err
	}
	return arts, nil
}
