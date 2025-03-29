package service

import (
	"context"
	"jiraAnalyzer/backend/internal/models"
	"jiraAnalyzer/backend/internal/repository"
)

type JiraClientService struct {
	repo *repository.Repository
}

func NewJiraClientService(repo *repository.Repository) *JiraClientService {
	return &JiraClientService{repo: repo}
}

func (c *JiraClientService) GetConnectorProjects(ctx context.Context, page, limit int, search string) ([]models.Project, models.PageInfo, error) {
	return c.repo.GetConnectorProjects(ctx, page, limit, search)
}

func (c *JiraClientService) UpdateConnectorProject(ctx context.Context, projectKeys []string) error {
	return c.repo.UpdateConnectorProject(ctx, projectKeys)
}
