package database

import (
	"database/sql"
	"fmt"
	"jiraAnalyzer/jiraConnector/internal/models"
)

func (r *JiraPostgres) SaveChangelogTx(tx *sql.Tx, changelogs []models.DBChangelog) error {
	stmt, err := tx.Prepare(`
        INSERT INTO status_changes (issue_id, created, from_status, to_status)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (issue_id, created) DO NOTHING
    `)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, cl := range changelogs {
		_, err := stmt.Exec(cl.IssueID, cl.Created, cl.FromStatus, cl.ToStatus)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}

	return nil
}
