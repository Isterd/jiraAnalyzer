package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"jiraAnalyzer/backend/internal/models"
)

type ProjectPostgres struct {
	db *sqlx.DB
}

func NewProjectPostgres(db *sqlx.DB) *ProjectPostgres {
	return &ProjectPostgres{db: db}
}

func (r *ProjectPostgres) CreateProject(ctx context.Context, project models.Project) (string, error) {
	var projectKey string
	query := "INSERT INTO projects (key, name, url) VALUES ($1, $2, $3) RETURNING key"
	err := r.db.QueryRowContext(ctx, query, project.Key, project.Name, project.URL).Scan(&projectKey)
	if err != nil {
		return "", fmt.Errorf("failed to create project: %w", err)
	}
	return projectKey, nil
}

func (r *ProjectPostgres) GetAllProjects(ctx context.Context, page, limit int, search string) ([]models.Project, models.PageInfo, error) {
	query := `
        SELECT key, name, url
        FROM projects
        WHERE ($1 = '' OR LOWER(name) LIKE LOWER($1) OR LOWER(key) LIKE LOWER($1))
        ORDER BY key
        LIMIT $2 OFFSET $3
    `

	var projects []models.Project
	offset := (page - 1) * limit
	err := r.db.SelectContext(ctx, &projects, query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, models.PageInfo{}, fmt.Errorf("failed to get projects: %w", err)
	}

	countQuery := `
        SELECT COUNT(*)
        FROM projects
        WHERE ($1 = '' OR LOWER(name) LIKE LOWER($1) OR LOWER(key) LIKE LOWER($1))
    `

	var totalCount int
	err = r.db.GetContext(ctx, &totalCount, countQuery, "%"+search+"%")
	if err != nil {
		return nil, models.PageInfo{}, fmt.Errorf("failed to count projects: %w", err)
	}

	pageInfo := models.PageInfo{
		CurrentPage: page,
		PageCount:   (totalCount + limit - 1) / limit,
		TotalCount:  totalCount,
	}

	return projects, pageInfo, nil
}

func (r *ProjectPostgres) GetProjectByID(ctx context.Context, id int) (models.Project, error) {
	var project models.Project
	query := "SELECT key, name, url FROM projects WHERE id = $1"
	err := r.db.GetContext(ctx, &project, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Project{}, fmt.Errorf("project with id %d not found", id)
		}
		return models.Project{}, fmt.Errorf("failed to get project: %w", err)
	}
	return project, nil
}

func (r *ProjectPostgres) UpdateProject(ctx context.Context, id int, project models.Project) error {
	query := "UPDATE projects SET key = $1, name = $2, url = $3 WHERE id = $4"
	_, err := r.db.ExecContext(ctx, query, project.Key, project.Name, project.URL, id)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	return nil
}

func (r *ProjectPostgres) DeleteProject(ctx context.Context, id int) error {
	query := "DELETE FROM projects WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}
