package memory

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

type OperationRepo struct{ s *Store }

func NewOperationRepo(s *Store) *OperationRepo { return &OperationRepo{s: s} }

var _ port.OperationRepo = (*OperationRepo)(nil)

func (r *OperationRepo) Create(_ context.Context, in model.Operation) (*model.Operation, error) {
	in.ID = uuid.New()
	in.CreatedAt = time.Now()
	r.s.mu.Lock()
	r.s.operations[in.ID] = in
	r.s.mu.Unlock()
	return &in, nil
}

func (r *OperationRepo) ListBySpec(_ context.Context, specID uuid.UUID) ([]model.Operation, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	var out []model.Operation
	for _, o := range r.s.operations {
		if o.SpecID == specID {
			out = append(out, o)
		}
	}
	return out, nil
}

func (r *OperationRepo) UpdateClassification(_ context.Context, id uuid.UUID, classification string, isDestructive bool) (*model.Operation, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	o, ok := r.s.operations[id]
	if !ok {
		return nil, ErrNotFound
	}
	o.Classification = classification
	o.IsDestructive = isDestructive
	r.s.operations[id] = o
	return &o, nil
}

func (r *OperationRepo) UpdateOverrides(_ context.Context, id uuid.UUID, overrides []byte) (*model.Operation, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	o, ok := r.s.operations[id]
	if !ok {
		return nil, ErrNotFound
	}
	o.Overrides = overrides
	r.s.operations[id] = o
	return &o, nil
}
