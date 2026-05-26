package memory

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

type SpecRepo struct{ s *Store }

func NewSpecRepo(s *Store) *SpecRepo { return &SpecRepo{s: s} }

var _ port.SpecRepo = (*SpecRepo)(nil)

func (r *SpecRepo) Create(_ context.Context, in model.Spec) (*model.Spec, error) {
	in.ID = uuid.New()
	in.CreatedAt = time.Now()
	r.s.mu.Lock()
	r.s.specs[in.ID] = in
	r.s.mu.Unlock()
	return &in, nil
}

func (r *SpecRepo) Get(_ context.Context, id uuid.UUID) (*model.Spec, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	s, ok := r.s.specs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &s, nil
}

func (r *SpecRepo) ListByProject(_ context.Context, projectID uuid.UUID) ([]model.Spec, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	var out []model.Spec
	for _, s := range r.s.specs {
		if s.ProjectID == projectID {
			out = append(out, s)
		}
	}
	return out, nil
}
