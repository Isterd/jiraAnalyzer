package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (r *JiraPostgres) GetOrCreateAuthor(displayName string) (int, error) {
	var authorID int

	createAuthorQuery := "INSERT INTO authors (display_name) VALUES ($1) ON CONFLICT (display_name) DO NOTHING RETURNING id"
	err := r.db.QueryRow(createAuthorQuery, displayName).Scan(&authorID)
	if errors.Is(err, sql.ErrNoRows) {
		err = r.db.QueryRow("SELECT id FROM authors WHERE display_name = $1", displayName).Scan(&authorID)
		if err != nil {
			return 0, fmt.Errorf("failed to get author ID: %w", err)
		}
	} else if err != nil {
		return 0, fmt.Errorf("failed to insert author: %w", err)
	}

	return authorID, nil
}
