package database

import (
	"database/sql"
	"fmt"
	"jiraAnalyzer/jiraConnector/internal/models"
)

func (r *JiraPostgres) SaveIssuesTx(tx *sql.Tx, issues []models.DBIssue) error {
	stmt, err := tx.Prepare(`
        INSERT INTO issues (
            key, project_key, created, updated, closed,
            summary, description, issue_type, priority, status,
            time_spent, creator_id, assignee_id
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
        ON CONFLICT (key) DO UPDATE SET
            updated = EXCLUDED.updated,
            status = EXCLUDED.status
    `)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, issue := range issues {
		_, err := stmt.Exec(
			issue.Key,
			issue.ProjectKey,
			issue.Created,
			issue.Updated,
			issue.Closed,
			issue.Summary,
			issue.Description,
			issue.Type,
			issue.Priority,
			issue.Status,
			issue.TimeSpent,
			issue.CreatorID,
			issue.AssigneeID,
		)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}

	return nil
}
