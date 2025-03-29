package service

import (
	"jiraAnalyzer/backend/internal/repository"
)

type Service struct {
	Projects   *ProjectService
	Issues     *IssueService
	Analytics  *AnalyticsService
	JiraClient *JiraClientService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Projects:   NewProjectService(repo),
		Issues:     NewIssueService(repo),
		Analytics:  NewAnalyticsService(repo),
		JiraClient: NewJiraClientService(repo),
	}
}
