package service

import (
	"context"
	"jiraAnalyzer/jiraConnector/internal/jiraclient"
	"jiraAnalyzer/jiraConnector/internal/models"
	"jiraAnalyzer/jiraConnector/internal/repository/database"
)

type JiraService struct {
	client *jiraclient.JiraClient
	repo   *database.Repository
}

func NewJiraService(client *jiraclient.JiraClient, repo *database.Repository) *JiraService {
	return &JiraService{
		client: client,
		repo:   repo,
	}
}

func (s *JiraService) GetAllProjects(ctx context.Context) ([]models.JiraProject, error) {
	return s.client.GetAllProjects(ctx)
}

func (s *JiraService) GetProjectIssues(ctx context.Context, projectKey string, startAt int) ([]models.JiraIssue, error) {
	return s.client.GetProjectIssues(ctx, projectKey, startAt)
}

func (s *JiraService) GetIssueCount(ctx context.Context, projectKey string) (int, error) {
	return s.client.GetIssueCount(ctx, projectKey)
}
