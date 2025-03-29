package repository

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"jiraAnalyzer/jiraConnector/internal/repository/database"
	"jiraAnalyzer/jiraConnector/internal/repository/jira"

	"jiraAnalyzer/jiraConnector/internal/models"
)

type JiraDB interface {
	// Проекты
	CheckProjectExists(ctx context.Context, projectKey string) (bool, error)
	SaveProject(ctx context.Context, project models.DBProject) error

	// Задачи
	SaveIssuesTx(tx *sql.Tx, issues []models.DBIssue) error
	SaveChangelogTx(tx *sql.Tx, changelogs []models.DBChangelog) error
	GetOrCreateAuthor(displayName string) (int, error)

	// Транзакции
	BeginTx(ctx context.Context) (*sql.Tx, error)
}

type JiraClient interface {
	// Клиент
	GetAllProjects(ctx context.Context) ([]models.JiraProject, error)
	GetProjectIssues(ctx context.Context, projectKey string, startAt int) ([]models.JiraIssue, error)
	GetIssueCount(ctx context.Context, projectKey string) (int, error)
}

type Repository struct {
	JiraDB
	JiraClient
}

func NewRepository(db *sqlx.DB, client *jira.Jira) *Repository {
	return &Repository{
		JiraDB:     database.NewJiraPostgres(db),
		JiraClient: client,
	}
}
