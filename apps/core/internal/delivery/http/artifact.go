package http

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/app"
)

// RegisterArtifacts registers run artifact endpoints.
func RegisterArtifacts(api huma.API, svc *app.ArtifactService) {
	huma.Get(api, "/api/runs/{id}/artifacts", func(ctx context.Context, in *IDParam) (*ArtifactListOutput, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, ErrBadRequest("artifacts", "invalid id")
		}
		as, err := svc.ListByRun(ctx, id)
		if err != nil {
			return nil, ErrInternal("artifacts", "failed to list artifacts", err)
		}
		return &ArtifactListOutput{Body: as}, nil
	})
}
