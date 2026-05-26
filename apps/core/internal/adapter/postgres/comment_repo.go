package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/fanboykun/smokery/apps/core/internal/adapter/postgres/db"
	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

type CommentRepo struct{ q *db.Queries }

func NewCommentRepo(pool *pgxpool.Pool) *CommentRepo {
	return &CommentRepo{q: db.New(pool)}
}

var _ port.CommentRepo = (*CommentRepo)(nil)

func (r *CommentRepo) Create(ctx context.Context, runID uuid.UUID, author, body string) (*model.Comment, error) {
	c, err := r.q.CreateComment(ctx, db.CreateCommentParams{
		RunID:  toPgUUID(runID),
		Author: author,
		Body:   body,
	})
	if err != nil {
		return nil, err
	}
	return commentToModel(c), nil
}

func (r *CommentRepo) ListByRun(ctx context.Context, runID uuid.UUID) ([]model.Comment, error) {
	cs, err := r.q.ListCommentsByRun(ctx, toPgUUID(runID))
	if err != nil {
		return nil, err
	}
	out := make([]model.Comment, len(cs))
	for i, c := range cs {
		out[i] = *commentToModel(c)
	}
	return out, nil
}

func commentToModel(c db.Comment) *model.Comment {
	return &model.Comment{
		ID:        fromPgUUID(c.ID),
		RunID:     fromPgUUID(c.RunID),
		Author:    c.Author,
		Body:      c.Body,
		CreatedAt: fromPgTimestamptz(c.CreatedAt),
	}
}
