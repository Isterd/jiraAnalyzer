package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"jiraAnalyzer/backend/internal/models"
	"jiraAnalyzer/backend/internal/repository/database"
	"jiraAnalyzer/backend/internal/repository/jira"
)

type Projects interface {
	GetAllProjects(ctx context.Context, page, limit int, search string) ([]models.Project, models.PageInfo, error)
	GetProjectByID(ctx context.Context, id int) (models.Project, error)
	CreateProject(ctx context.Context, project models.Project) (string, error)
	UpdateProject(ctx context.Context, id int, project models.Project) error
	DeleteProject(ctx context.Context, id int) error
}

type Issues interface {
	GetAllIssues(ctx context.Context, page, limit int) ([]models.Issue, error)
	GetIssueById(ctx context.Context, id int) (models.Issue, error)
	CreateIssue(ctx context.Context, issue models.Issue) error
	UpdateIssue(ctx context.Context, issue models.Issue) error
	DeleteIssue(ctx context.Context, id int) error
	AddStatusChange(ctx context.Context, change models.StatusChange) error
}

type Authors interface {
	GetAuthorById(ctx context.Context, id int) (models.Author, error)
	CreateAuthor(ctx context.Context, author models.Author) error
}

type Analytics interface {
	GetProjectAnalytics(ctx context.Context, projectKey string) (models.ProjectAnalytics, error)
	GetAnalytics(ctx context.Context, projectKey string, taskNumber int) ([]byte, error)
	IsProjectAnalyzed(ctx context.Context, projectKey string) (bool, error)
	DeleteProjectAnalytics(ctx context.Context, projectKey string) error
	SaveAnalytics(ctx context.Context, projectKey string, taskNumber int, data []byte) error
	CalculateOpenTimeHistogram(ctx context.Context, projectKey string) ([]models.HistogramData, error)
	CalculateStatusTimeDistribution(ctx context.Context, projectKey string) ([]models.StatusTimeData, error)
	CalculateActivityGraph(ctx context.Context, projectKey string) ([]models.ActivityData, error)
	CalculateComplexityGraph(ctx context.Context, projectKey string) ([]models.ComplexityData, error)
	CalculatePriorityDistribution(ctx context.Context, projectKey string) ([]models.PriorityData, error)
	CalculatePriorityDistributionClosedTasks(ctx context.Context, projectKey string) ([]models.PriorityData, error)
}

type JiraClient interface {
	GetConnectorProjects(ctx context.Context, page, limit int, search string) ([]models.Project, models.PageInfo, error)
	UpdateConnectorProject(ctx context.Context, projectKeys []string) error
}

type Repository struct {
	Projects
	Issues
	Authors
	Analytics
	JiraClient
}

func NewRepository(db *sqlx.DB, url string) *Repository {
	return &Repository{
		Projects:   database.NewProjectPostgres(db),
		Issues:     database.NewIssuePostgres(db),
		Authors:    database.NewAuthorPostgres(db),
		Analytics:  database.NewAnalyticsPostgres(db),
		JiraClient: jira.NewHTTPJiraClient(url),
	}
}
