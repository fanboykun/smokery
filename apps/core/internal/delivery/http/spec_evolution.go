package http

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/app"
)

type SpecDiffInput struct {
	FromID string `path:"from-id" format:"uuid" doc:"Source spec UUID"`
	ToID   string `path:"to-id" format:"uuid" doc:"Target spec UUID"`
}

type SpecDiffOutput struct {
	Body app.SpecDiff
}

type ImpactInput struct {
	ProjectIDParam
	SpecID string `path:"spec-id" format:"uuid" doc:"Spec UUID"`
}

type ImpactOutput struct {
	Body app.ImpactAnalysis
}

func RegisterSpecEvolution(api huma.API, svc *app.SpecEvolutionService) {
	huma.Get(api, "/api/specs/{from-id}/diff/{to-id}", func(ctx context.Context, in *SpecDiffInput) (*SpecDiffOutput, error) {
		fromID, err := uuid.Parse(in.FromID)
		if err != nil {
			return nil, ErrBadRequest("spec-diff", "invalid from-id")
		}
		toID, err := uuid.Parse(in.ToID)
		if err != nil {
			return nil, ErrBadRequest("spec-diff", "invalid to-id")
		}
		diff, err := svc.Diff(ctx, fromID, toID)
		if err != nil {
			return nil, ErrInternal("spec-diff", "failed to compute diff", err)
		}
		return &SpecDiffOutput{Body: *diff}, nil
	})

	huma.Get(api, "/api/projects/{project-id}/impact/spec/{spec-id}", func(ctx context.Context, in *ImpactInput) (*ImpactOutput, error) {
		projectID, err := uuid.Parse(in.ProjectID)
		if err != nil {
			return nil, ErrBadRequest("impact", "invalid project id")
		}
		specID, err := uuid.Parse(in.SpecID)
		if err != nil {
			return nil, ErrBadRequest("impact", "invalid spec id")
		}
		analysis, err := svc.Impact(ctx, projectID, specID)
		if err != nil {
			return nil, ErrInternal("impact", "failed to analyze impact", err)
		}
		return &ImpactOutput{Body: *analysis}, nil
	})
}
