package http

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/app"
)

// RegisterProjects registers project CRUD endpoints.
func RegisterProjects(api huma.API, svc *app.ProjectService) {
	huma.Post(api, "/api/projects", func(ctx context.Context, in *CreateProjectInput) (*ProjectOutput, error) {
		p, err := svc.Create(ctx, in.Body.Name, in.Body.Description)
		if err != nil {
			return nil, ErrInternal("projects", "failed to create project", err)
		}
		return &ProjectOutput{Body: *p}, nil
	})

	huma.Get(api, "/api/projects", func(ctx context.Context, in *struct{}) (*ProjectListOutput, error) {
		ps, err := svc.List(ctx)
		if err != nil {
			return nil, ErrInternal("projects", "failed to list projects", err)
		}
		return &ProjectListOutput{Body: ps}, nil
	})

	huma.Get(api, "/api/projects/{id}", func(ctx context.Context, in *IDParam) (*ProjectOutput, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, ErrBadRequest("projects/{id}", "invalid project id")
		}
		p, err := svc.Get(ctx, id)
		if err != nil {
			return nil, ErrNotFound("projects/{id}", "project not found")
		}
		return &ProjectOutput{Body: *p}, nil
	})

	huma.Put(api, "/api/projects/{id}", func(ctx context.Context, in *UpdateProjectInput) (*ProjectOutput, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, ErrBadRequest("projects/{id}", "invalid project id")
		}
		p, err := svc.Update(ctx, id, in.Body.Name, in.Body.Description)
		if err != nil {
			return nil, ErrInternal("projects/{id}", "failed to update project", err)
		}
		return &ProjectOutput{Body: *p}, nil
	})

	huma.Delete(api, "/api/projects/{id}", func(ctx context.Context, in *IDParam) (*struct{}, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, ErrBadRequest("projects/{id}", "invalid project id")
		}
		if err := svc.Delete(ctx, id); err != nil {
			return nil, ErrInternal("projects/{id}", "failed to delete project", err)
		}
		return nil, nil
	})
}
