package controller

import (
	"github.com/sirupsen/logrus"
	"jiraAnalyzer/backend/internal/config"
	"jiraAnalyzer/backend/internal/service"
)

type Controller struct {
	*ProjectController
	*IssueController
	*AnalyticsController
	*JiraController
}

func NewController(service *service.Service, logger *logrus.Logger, cfg config.Backend) *Controller {
	return &Controller{
		ProjectController:   NewProjectController(service.Projects),
		IssueController:     NewIssueController(service.Issues),
		AnalyticsController: NewAnalyticsController(service.Analytics, cfg),
		JiraController:      NewJiraController(service.JiraClient),
	}
}
