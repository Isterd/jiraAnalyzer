package database

import (
	"context"
	"jiraAnalyzer/jiraConnector/internal/models"
)

func (r *JiraPostgres) CheckProjectExists(ctx context.Context, projectKey string) (bool, error) {
	var exists bool
	err := r.db.QueryRowxContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM projects WHERE key = $1)",
		projectKey).Scan(&exists) // <-- QueryRowxContext
	return exists, err
}

func (r *JiraPostgres) SaveProject(ctx context.Context, project models.DBProject) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		"INSERT INTO projects (key, name, url) VALUES ($1, $2, $3) ON CONFLICT (key) DO NOTHING",
		project.Key, project.Name, project.URL,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}
