package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/fanboykun/smokery/apps/core/internal/adapter/postgres/db"
	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

type SpecRepo struct{ q *db.Queries }

func NewSpecRepo(pool *pgxpool.Pool) *SpecRepo {
	return &SpecRepo{q: db.New(pool)}
}

var _ port.SpecRepo = (*SpecRepo)(nil)

func (r *SpecRepo) Create(ctx context.Context, in model.Spec) (*model.Spec, error) {
	s, err := r.q.CreateSpec(ctx, db.CreateSpecParams{
		ProjectID: toPgUUID(in.ProjectID),
		Version:   in.Version,
		Title:     in.Title,
		Raw:       in.Raw,
		Analysis:  in.Analysis,
	})
	if err != nil {
		return nil, err
	}
	return specToModel(s), nil
}

func (r *SpecRepo) Get(ctx context.Context, id uuid.UUID) (*model.Spec, error) {
	s, err := r.q.GetSpec(ctx, toPgUUID(id))
	if err != nil {
		return nil, err
	}
	return specToModel(s), nil
}

func (r *SpecRepo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]model.Spec, error) {
	ss, err := r.q.ListSpecsByProject(ctx, toPgUUID(projectID))
	if err != nil {
		return nil, err
	}
	out := make([]model.Spec, len(ss))
	for i, s := range ss {
		out[i] = *specToModel(s)
	}
	return out, nil
}

func specToModel(s db.Spec) *model.Spec {
	return &model.Spec{
		ID:        fromPgUUID(s.ID),
		ProjectID: fromPgUUID(s.ProjectID),
		Version:   s.Version,
		Title:     s.Title,
		Raw:       s.Raw,
		Analysis:  s.Analysis,
		CreatedAt: fromPgTimestamptz(s.CreatedAt),
	}
}
