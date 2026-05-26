package app

import (
	"context"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

type CommentService struct {
	comments port.CommentRepo
}

func NewCommentService(c port.CommentRepo) *CommentService {
	return &CommentService{comments: c}
}

func (s *CommentService) Create(ctx context.Context, runID uuid.UUID, author, body string) (*model.Comment, error) {
	return s.comments.Create(ctx, runID, author, body)
}

func (s *CommentService) ListByRun(ctx context.Context, runID uuid.UUID) ([]model.Comment, error) {
	return s.comments.ListByRun(ctx, runID)
}
