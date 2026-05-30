package app

import (
	"context"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

type ProjectService struct {
	projects port.ProjectRepo
	specs    port.SpecRepo
	runs     port.RunRepo
}

func NewProjectService(p port.ProjectRepo, deps ...any) *ProjectService {
	svc := &ProjectService{projects: p}
	for _, dep := range deps {
		switch repo := dep.(type) {
		case port.SpecRepo:
			svc.specs = repo
		case port.RunRepo:
			svc.runs = repo
		}
	}
	return svc
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

func (s *ProjectService) ListWithStats(ctx context.Context) ([]model.ProjectWithStats, error) {
	projects, err := s.projects.List(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]model.ProjectWithStats, 0, len(projects))
	for _, project := range projects {
		specCount := 0
		if s.specs != nil {
			specs, err := s.specs.ListByProject(ctx, project.ID)
			if err != nil {
				return nil, err
			}
			specCount = len(specs)
		}

		stats := model.ProjectHealthStats{}
		var lastRun *model.LastRunInfo
		if s.runs != nil {
			runs, err := s.runs.ListByProject(ctx, project.ID)
			if err != nil {
				return nil, err
			}

			stats.TotalRuns = len(runs)
			for _, run := range runs {
				switch run.Status {
				case "passed", "completed":
					stats.PassedRuns++
				case "failed", "error":
					stats.FailedRuns++
				}

				if lastRun == nil || run.CreatedAt.After(lastRun.CreatedAt) {
					duration := int64(0)
					if run.StartedAt != nil && run.FinishedAt != nil {
						duration = run.FinishedAt.Sub(*run.StartedAt).Milliseconds()
					}
					lastRun = &model.LastRunInfo{
						ID:        run.ID,
						Status:    run.Status,
						CreatedAt: run.CreatedAt,
						Duration:  duration,
					}
				}
			}
		}

		if stats.TotalRuns > 0 {
			stats.HealthPercentage = float64(stats.PassedRuns) / float64(stats.TotalRuns) * 100
		}

		out = append(out, model.ProjectWithStats{
			ID:          project.ID,
			Name:        project.Name,
			Description: project.Description,
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt,
			SpecCount:   specCount,
			LastRun:     lastRun,
			Stats:       stats,
		})
	}

	return out, nil
}

func (s *ProjectService) Update(ctx context.Context, id uuid.UUID, name, description string) (*model.Project, error) {
	return s.projects.Update(ctx, id, name, description)
}

func (s *ProjectService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.projects.Delete(ctx, id)
}
