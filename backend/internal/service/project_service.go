package service

import (
	"context"
	"fmt"
	"jiraAnalyzer/backend/internal/models"
	"jiraAnalyzer/backend/internal/repository"
)

type ProjectService struct {
	repo *repository.Repository
}

func NewProjectService(repo *repository.Repository) *ProjectService {
	return &ProjectService{
		repo: repo,
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, project models.Project) (string, error) {
	return s.repo.CreateProject(ctx, project)
}

func (s *ProjectService) GetAllProjects(ctx context.Context, page, limit int, search string) ([]models.Project, models.PageInfo, error) {
	return s.repo.GetAllProjects(ctx, page, limit, search)
}

func (s *ProjectService) GetProjectByID(ctx context.Context, id int) (models.Project, error) {
	project, err := s.repo.GetProjectByID(ctx, id)
	if err != nil {
		return models.Project{}, fmt.Errorf("error while getting project by id: %w", err)
	}

	return project, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, id int, project models.Project) error {
	return s.repo.UpdateProject(ctx, id, project)
}

func (s *ProjectService) DeleteProject(ctx context.Context, id int) error {
	return s.repo.DeleteProject(ctx, id)
}
