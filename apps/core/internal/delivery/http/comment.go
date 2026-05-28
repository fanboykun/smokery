package http

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/app"
)

// RegisterComments registers run comment endpoints.
func RegisterComments(api huma.API, svc *app.CommentService) {
	huma.Post(api, "/api/runs/{id}/comments", func(ctx context.Context, in *CreateCommentInput) (*CommentOutput, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, ErrBadRequest("comments", "invalid id")
		}
		c, err := svc.Create(ctx, id, in.Body.Author, in.Body.Body)
		if err != nil {
			return nil, ErrInternal("comments", "failed to create comment", err)
		}
		return &CommentOutput{Body: *c}, nil
	})

	huma.Get(api, "/api/runs/{id}/comments", func(ctx context.Context, in *IDParam) (*CommentListOutput, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, ErrBadRequest("comments", "invalid id")
		}
		cs, err := svc.ListByRun(ctx, id)
		if err != nil {
			return nil, ErrInternal("comments", "failed to list comments", err)
		}
		return &CommentListOutput{Body: cs}, nil
	})
}
