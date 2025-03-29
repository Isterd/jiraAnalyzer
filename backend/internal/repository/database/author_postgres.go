package database

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"jiraAnalyzer/backend/internal/models"
)

type AuthorPostgres struct {
	db *sqlx.DB
}

func NewAuthorPostgres(db *sqlx.DB) *AuthorPostgres {
	return &AuthorPostgres{db: db}
}

func (r *AuthorPostgres) GetAuthorById(ctx context.Context, id int) (models.Author, error) {
	var author models.Author
	query := "SELECT id, display_name FROM authors WHERE id = $1"
	err := r.db.GetContext(ctx, &author, query, id)
	if err != nil {
		return models.Author{}, fmt.Errorf("failed to get author by ID: %w", err)
	}
	return author, nil
}

func (r *AuthorPostgres) CreateAuthor(ctx context.Context, author models.Author) error {
	query := `
        INSERT INTO authors (display_name)
        VALUES ($1)
        RETURNING id
    `
	err := r.db.QueryRowContext(ctx, query, author.DisplayName).Scan(&author.ID)
	return err
}
