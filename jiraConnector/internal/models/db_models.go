package models

import "time"

type PageInfo struct {
	CurrentPage int `json:"currentPage"`
	PageCount   int `json:"pageCount"`
	TotalCount  int `json:"totalCount"`
}

type DBProject struct {
	Key       string `db:"key"`
	Name      string `db:"name"`
	URL       string `db:"url"`
	CreatedAt string `db:"created_at"`
}

type DBIssue struct {
	JiraID         string     `db:"jira_id"`
	ProjectKey     string     `db:"project_key"`
	Key            string     `db:"key"`
	Created        time.Time  `db:"created"`
	Updated        time.Time  `db:"updated"`
	ResolutionDate *time.Time `db:"resolution_date"`
	Summary        string     `db:"summary"`
	Description    string     `db:"description"`
	Type           string     `db:"issue_type"`
	Priority       string     `db:"priority"`
	Status         string     `db:"status"`
	TimeSpent      int        `db:"time_spent"`
	CreatorID      int        `db:"creator_id"`
	AssigneeID     *int       `db:"assignee_id"`
}

type DBChangelog struct {
	IssueID    string    `db:"issue_id"`
	AuthorID   int       `db:"author_id"`
	Created    time.Time `db:"created"`
	FromStatus string    `db:"from_status"`
	ToStatus   string    `db:"to_status"`
}

type DBAuthor struct {
	ID          int    `db:"id"`
	DisplayName string `db:"display_name"`
}
