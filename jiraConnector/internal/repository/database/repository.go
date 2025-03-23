package database

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"

	"jiraAnalyzer/jiraConnector/internal/models"
)

type Jira interface {
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

type Repository struct {
	Jira
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Jira: NewJiraPostgres(db),
	}
}
