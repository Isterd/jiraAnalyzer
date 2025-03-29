package models

import "time"

type PageInfo struct {
	CurrentPage int `json:"currentPage"`
	PageCount   int `json:"pageCount"`
	TotalCount  int `json:"totalCount"`
}

type DBProject struct {
	ID        int    `db:"id"`
	Key       string `db:"key"`
	Name      string `db:"name"`
	URL       string `db:"url"`
	CreatedAt string `db:"created_at"`
}

type DBIssue struct {
	ID             int        `db:"id"`
	Key            string     `db:"key"`
	ProjectKey     string     `db:"project_key"`
	Created        time.Time  `db:"created"`
	Updated        time.Time  `db:"updated"`
	Closed         *time.Time `db:"closed"`
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
	ID         int       `db:"id"`
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
