package app

import (
	"context"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

type ProjectService struct {
	projects port.ProjectRepo
}

func NewProjectService(p port.ProjectRepo) *ProjectService {
	return &ProjectService{projects: p}
}

func (s *ProjectService) Create(ctx context.Context, name, description string) (*model.Project, error) {
	return s.projects.Create(ctx, name, description)
}

func (s *ProjectService) Get(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	return s.projects.Get(ctx, id)
}

func (s *ProjectService) List(ctx context.Context) ([]model.Project, error) {
	return s.projects.List(ctx)
}

func (s *ProjectService) Update(ctx context.Context, id uuid.UUID, name, description string) (*model.Project, error) {
	return s.projects.Update(ctx, id, name, description)
}

func (s *ProjectService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.projects.Delete(ctx, id)
}
