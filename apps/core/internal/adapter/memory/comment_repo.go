package memory

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

type CommentRepo struct{ s *Store }

func NewCommentRepo(s *Store) *CommentRepo { return &CommentRepo{s: s} }

var _ port.CommentRepo = (*CommentRepo)(nil)

func (r *CommentRepo) Create(_ context.Context, runID uuid.UUID, author, body string) (*model.Comment, error) {
	c := model.Comment{ID: uuid.New(), RunID: runID, Author: author, Body: body, CreatedAt: time.Now()}
	r.s.mu.Lock()
	r.s.comments[c.ID] = c
	r.s.mu.Unlock()
	return &c, nil
}

func (r *CommentRepo) ListByRun(_ context.Context, runID uuid.UUID) ([]model.Comment, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	var out []model.Comment
	for _, c := range r.s.comments {
		if c.RunID == runID {
			out = append(out, c)
		}
	}
	return out, nil
}
