package memory

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

type ArtifactRepo struct{ s *Store }

func NewArtifactRepo(s *Store) *ArtifactRepo { return &ArtifactRepo{s: s} }

var _ port.ArtifactRepo = (*ArtifactRepo)(nil)

func (r *ArtifactRepo) Create(_ context.Context, runID uuid.UUID, typ, path string) (*model.Artifact, error) {
	a := model.Artifact{ID: uuid.New(), RunID: runID, Type: typ, Path: path, CreatedAt: time.Now()}
	r.s.mu.Lock()
	r.s.artifacts[a.ID] = a
	r.s.mu.Unlock()
	return &a, nil
}

func (r *ArtifactRepo) ListByRun(_ context.Context, runID uuid.UUID) ([]model.Artifact, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	var out []model.Artifact
	for _, a := range r.s.artifacts {
		if a.RunID == runID {
			out = append(out, a)
		}
	}
	return out, nil
}

func (r *ArtifactRepo) DeleteByRun(_ context.Context, runID uuid.UUID) ([]model.Artifact, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	var deleted []model.Artifact
	for id, a := range r.s.artifacts {
		if a.RunID == runID {
			deleted = append(deleted, a)
			delete(r.s.artifacts, id)
		}
	}
	return deleted, nil
}
