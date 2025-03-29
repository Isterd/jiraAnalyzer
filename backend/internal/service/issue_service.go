package service

import (
	"context"
	"jiraAnalyzer/backend/internal/models"
	"jiraAnalyzer/backend/internal/repository"
)

type IssueService struct {
	repo *repository.Repository
}

func NewIssueService(repo *repository.Repository) *IssueService {
	return &IssueService{repo: repo}
}

func (s *IssueService) AddStatusChange(ctx context.Context, change models.StatusChange) error {
	return s.repo.AddStatusChange(ctx, change)
}

func (s *IssueService) GetAllIssues(ctx context.Context, page, limit int) ([]models.Issue, error) {
	return s.repo.GetAllIssues(ctx, page, limit)
}

func (s *IssueService) GetIssueById(ctx context.Context, id int) (models.Issue, error) {
	return s.repo.GetIssueById(ctx, id)
}

func (s *IssueService) CreateIssue(ctx context.Context, issue models.Issue) error {
	return s.repo.CreateIssue(ctx, issue)
}

func (s *IssueService) UpdateIssue(ctx context.Context, issue models.Issue) error {
	return s.repo.UpdateIssue(ctx, issue)
}

func (s *IssueService) DeleteIssue(ctx context.Context, id int) error {
	return s.repo.DeleteIssue(ctx, id)
}
