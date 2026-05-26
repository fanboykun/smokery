package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/fanboykun/smokery/apps/core/internal/adapter/postgres/db"
	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

type ProjectRepo struct{ q *db.Queries }

func NewProjectRepo(pool *pgxpool.Pool) *ProjectRepo {
	return &ProjectRepo{q: db.New(pool)}
}

var _ port.ProjectRepo = (*ProjectRepo)(nil)

func (r *ProjectRepo) Create(ctx context.Context, name, description string) (*model.Project, error) {
	p, err := r.q.CreateProject(ctx, db.CreateProjectParams{Name: name, Description: description})
	if err != nil {
		return nil, err
	}
	return projectToModel(p), nil
}

func (r *ProjectRepo) Get(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	p, err := r.q.GetProject(ctx, toPgUUID(id))
	if err != nil {
		return nil, err
	}
	return projectToModel(p), nil
}

func (r *ProjectRepo) List(ctx context.Context) ([]model.Project, error) {
	ps, err := r.q.ListProjects(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]model.Project, len(ps))
	for i, p := range ps {
		out[i] = *projectToModel(p)
	}
	return out, nil
}

func (r *ProjectRepo) Update(ctx context.Context, id uuid.UUID, name, description string) (*model.Project, error) {
	p, err := r.q.UpdateProject(ctx, db.UpdateProjectParams{ID: toPgUUID(id), Name: name, Description: description})
	if err != nil {
		return nil, err
	}
	return projectToModel(p), nil
}

func (r *ProjectRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteProject(ctx, toPgUUID(id))
}

func projectToModel(p db.Project) *model.Project {
	return &model.Project{
		ID:          fromPgUUID(p.ID),
		Name:        p.Name,
		Description: p.Description,
		CreatedAt:   fromPgTimestamptz(p.CreatedAt),
		UpdatedAt:   fromPgTimestamptz(p.UpdatedAt),
	}
}
