package memory

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

var ErrNotFound = errors.New("not found")

type ProjectRepo struct{ s *Store }

func NewProjectRepo(s *Store) *ProjectRepo { return &ProjectRepo{s: s} }

var _ port.ProjectRepo = (*ProjectRepo)(nil)

func (r *ProjectRepo) Create(_ context.Context, name, description string) (*model.Project, error) {
	now := time.Now()
	p := model.Project{ID: uuid.New(), Name: name, Description: description, CreatedAt: now, UpdatedAt: now}
	r.s.mu.Lock()
	r.s.projects[p.ID] = p
	r.s.mu.Unlock()
	return &p, nil
}

func (r *ProjectRepo) Get(_ context.Context, id uuid.UUID) (*model.Project, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	p, ok := r.s.projects[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &p, nil
}

func (r *ProjectRepo) List(_ context.Context) ([]model.Project, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	out := make([]model.Project, 0, len(r.s.projects))
	for _, p := range r.s.projects {
		out = append(out, p)
	}
	return out, nil
}

func (r *ProjectRepo) Update(_ context.Context, id uuid.UUID, name, description string) (*model.Project, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	p, ok := r.s.projects[id]
	if !ok {
		return nil, ErrNotFound
	}
	p.Name = name
	p.Description = description
	p.UpdatedAt = time.Now()
	r.s.projects[id] = p
	return &p, nil
}

func (r *ProjectRepo) Delete(_ context.Context, id uuid.UUID) error {
	r.s.mu.Lock()
	delete(r.s.projects, id)
	r.s.mu.Unlock()
	return nil
}
