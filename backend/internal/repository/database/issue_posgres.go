package database

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"jiraAnalyzer/backend/internal/models"
	"log"
)

type IssuePostgres struct {
	db *sqlx.DB
}

func NewIssuePostgres(db *sqlx.DB) *IssuePostgres {
	return &IssuePostgres{db: db}
}

func (r *IssuePostgres) AddStatusChange(ctx context.Context, change models.StatusChange) error {
	query := `
        INSERT INTO status_changes (issue_id, author_id, created, from_status, to_status)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err := r.db.ExecContext(ctx, query, change.IssueID, change.AuthorID, change.Created, change.FromStatus, change.ToStatus)
	return err
}

// GetAllIssues получает все задачи с пагинацией
func (r *IssuePostgres) GetAllIssues(ctx context.Context, page, limit int) ([]models.Issue, error) {
	var issues []models.Issue
	offset := (page - 1) * limit
	query := `
        SELECT key, project_key, created, updated, closed, summary, description, issue_type, priority, status, time_spent, creator_id, assignee_id
        FROM issues
        LIMIT $1 OFFSET $2
    `
	err := r.db.SelectContext(ctx, &issues, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get issues from db: %w", err)
	}
	return issues, nil
}

// GetIssueById получает задачу по её ID
func (r *IssuePostgres) GetIssueById(ctx context.Context, id int) (models.Issue, error) {
	var issue models.Issue
	query := `
        SELECT key, project_key, created, updated, summary, description, issue_type, priority, status, time_spent, creator_id, assignee_id
        FROM issues
        WHERE id = $1
    `
	err := r.db.GetContext(ctx, &issue, query, id)
	if err != nil {
		return models.Issue{}, fmt.Errorf("failed to get issue by ID: %w", err)
	}
	return issue, nil
}

// CreateIssue создает новую задачу
func (r *IssuePostgres) CreateIssue(ctx context.Context, issue models.Issue) error {
	if err := r.validateIssue(ctx, issue); err != nil {
		return fmt.Errorf("failed to validate issue: %w", err)
	}

	query := `
        INSERT INTO issues (key, project_key, created, updated, summary, description, issue_type, priority, status, time_spent, creator_id, assignee_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    `
	_, err := r.db.ExecContext(ctx, query,
		issue.Key, issue.ProjectKey, issue.Created, issue.Updated,
		issue.Summary, issue.Description, issue.IssueType, issue.Priority, issue.Status, issue.TimeSpent,
		issue.CreatorID, issue.AssigneeID)
	if err != nil {
		return fmt.Errorf("failed to create issue: %w", err)
	}

	return nil
}

// UpdateIssue обновляет существующую задачу
func (r *IssuePostgres) UpdateIssue(ctx context.Context, issue models.Issue) error {
	if err := r.validateIssue(ctx, issue); err != nil {
		return fmt.Errorf("failed to validate issue: %w", err)
	}

	query := `
        UPDATE issues
        SET key = $1, project_key = $2, created = $3, updated = $4, summary = $5, description = $6, issue_type = $7, priority = $8, status = $9, time_spent = $10, creator_id = $11, assignee_id = $12
        WHERE id = $13
    `
	_, err := r.db.ExecContext(ctx, query,
		issue.Key, issue.ProjectKey, issue.Created, issue.Updated,
		issue.Summary, issue.Description, issue.IssueType, issue.Priority, issue.Status, issue.TimeSpent,
		issue.CreatorID, issue.AssigneeID, issue.ID)
	if err != nil {
		return fmt.Errorf("failed to update issue: %w", err)
	}

	return nil
}

// DeleteIssue удаляет задачу по её ID
func (r *IssuePostgres) DeleteIssue(ctx context.Context, id int) error {
	query := "DELETE FROM issues WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete issue: %w", err)
	}

	return nil
}

func (r *IssuePostgres) validateIssue(ctx context.Context, issue models.Issue) error {
	projectQuery := `SELECT EXISTS(SELECT 1 FROM projects WHERE key = $1)`
	var projectExists bool
	err := r.db.QueryRowContext(ctx, projectQuery, issue.ProjectKey).Scan(&projectExists)
	if err != nil {
		return fmt.Errorf("failed to check project existence: %w", err)
	}
	if !projectExists {
		log.Printf("Project with key %s does not exist in the database", issue.ProjectKey)
		return fmt.Errorf("project with key %s does not exist", issue.ProjectKey)
	}

	authorQuery := `SELECT EXISTS(SELECT 1 FROM authors WHERE id = $1)`
	var authorExists bool
	err = r.db.QueryRowContext(ctx, authorQuery, issue.CreatorID).Scan(&authorExists)
	if err != nil {
		return fmt.Errorf("failed to check creator existence: %w", err)
	}
	if !authorExists {
		log.Printf("Creator with ID %d does not exist in the database", issue.CreatorID)
		return fmt.Errorf("creator with ID %d does not exist", issue.CreatorID)
	}

	// Проверяем существование исполнителя задачи, если он указан
	if issue.AssigneeID != nil {
		assigneeQuery := `SELECT EXISTS(SELECT 1 FROM authors WHERE id = $1)`
		var assigneeExists bool
		err = r.db.QueryRowContext(ctx, assigneeQuery, *issue.AssigneeID).Scan(&assigneeExists)
		if err != nil {
			return fmt.Errorf("failed to check assignee existence: %w", err)
		}
		if !assigneeExists {
			log.Printf("Assignee with ID %d does not exist in the database", *issue.AssigneeID)
			return fmt.Errorf("assignee with ID %d does not exist", *issue.AssigneeID)
		}
	}

	return nil
}
