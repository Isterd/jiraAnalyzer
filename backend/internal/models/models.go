package models

import "time"

type Project struct {
	Key       string    `json:"key" db:"key"`
	Name      string    `json:"name" db:"name"`
	URL       string    `json:"url" db:"url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type PageInfo struct {
	CurrentPage int `json:"currentPage"`
	PageCount   int `json:"pageCount"`
	TotalCount  int `json:"totalCount"`
}

type Issue struct {
	ID          int        `json:"id" db:"id"`
	Key         string     `json:"key" db:"key"`
	ProjectKey  string     `json:"project_key" db:"project_key"`
	Created     time.Time  `json:"created" db:"created"`
	Updated     time.Time  `json:"updated" db:"updated"`
	Closed      *time.Time `json:"closed,omitempty" db:"closed"`
	Summary     string     `json:"summary" db:"summary"`
	Description string     `json:"description" db:"description"`
	IssueType   string     `json:"issue_type" db:"issue_type"`
	Priority    string     `json:"priority" db:"priority"`
	Status      string     `json:"status" db:"status"`
	TimeSpent   int        `json:"time_spent" db:"time_spent"`
	CreatorID   int        `json:"creator_id" db:"creator_id"`
	AssigneeID  *int       `json:"assignee_id,omitempty" db:"assignee_id"` // Может быть NULL
}

type StatusChange struct {
	ID         int       `json:"id" db:"id"`
	IssueID    string    `json:"issue_id" db:"issue_id"`
	AuthorID   int       `json:"author_id" db:"author_id"`
	Created    time.Time `json:"created" db:"created"`
	FromStatus string    `json:"from_status" db:"from_status"`
	ToStatus   string    `json:"to_status" db:"to_status"`
}

type Author struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type ProjectAnalytics struct {
	TotalIssues        int     `json:"total_issues"`
	ClosedIssues       int     `json:"closed_issues"`
	OpenIssues         int     `json:"open_issues"`
	ReopenIssues       int     `json:"reopen_issues"`
	ResolvedIssues     int     `json:"resolved_issues"`
	InProgressIssues   int     `json:"in_progress_issues"`
	AverageTimeIssues  float64 `json:"average_time_issues"`
	AverageCountIssues float64 `json:"average_count_issues"`
}

type HistogramData struct {
	DayInterval int `db:"day_interval"`
	TaskCount   int `db:"task_count"`
}

// StatusTimeData представляет данные о распределении времени задач по состояниям.
type StatusTimeData struct {
	Status      string `db:"status"`
	DayInterval int    `db:"day_interval"`
	TaskCount   int    `db:"task_count"`
}

type ActivityData struct {
	Day              string `db:"day"`
	CumulativeOpened int    `db:"cumulative_opened"`
	CumulativeClosed int    `db:"cumulative_closed"`
}

type ComplexityData struct {
	TimeSpentHours float64 `db:"time_spent_hours"`
	TaskCount      int     `db:"task_count"`
}

// PriorityData представляет данные о распределении задач по приоритетам.
type PriorityData struct {
	Priority  string `db:"priority"`   // Название приоритета
	TaskCount int    `db:"task_count"` // Количество задач
}

type ComparisonProjects struct {
	TotalTasks        int     `db:"total_tasks"`
	OpenTasks         int     `db:"open_tasks"`
	ClosedTasks       int     `db:"closed_tasks"`
	AverageTimeIssues float64 `db:"average_time_issues"`
}
