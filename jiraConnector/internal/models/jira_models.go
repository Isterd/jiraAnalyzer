package models

type JiraProject struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	URL  string `json:"self"`
}

// JiraSearchResponse представляет структуру ответа JiraDB API для эндпоинта /rest/api/2/search.
type JiraSearchResponse struct {
	Expand     string      `json:"expand"`
	StartAt    int         `json:"startAt"`
	MaxResults int         `json:"maxResults"`
	Total      int         `json:"total"`
	Issues     []JiraIssue `json:"issues"`
}

type JiraIssue struct {
	Key       string        `json:"key"`
	Fields    JiraFields    `json:"fields"`
	Changelog JiraChangelog `json:"changelog"`
}

type JiraFields struct {
	Created     string       `json:"created"`
	Updated     string       `json:"updated"`
	Summary     string       `json:"summary"`
	Description string       `json:"description"`
	IssueType   JiraType     `json:"issuetype"`
	Priority    JiraPriority `json:"priority"`
	Status      JiraStatus   `json:"status"`
	TimeSpent   int          `json:"timespent"`
	Creator     JiraAuthor   `json:"creator"`
	Assignee    *JiraAuthor  `json:"assignee"`
}

type JiraResolution struct {
	Date string `json:"date"`
}

type JiraType struct {
	Name string `json:"name"`
}

type JiraPriority struct {
	Name string `json:"name"`
}

type JiraStatus struct {
	Name string `json:"name"`
}

type JiraAuthor struct {
	DisplayName string `json:"displayName"`
}

type JiraChangelog struct {
	Histories []JiraHistory `json:"histories"`
}

type JiraHistory struct {
	Created string            `json:"created"`
	Author  JiraAuthor        `json:"author"`
	Items   []JiraHistoryItem `json:"items"`
}

type JiraHistoryItem struct {
	Field      string `json:"field"`
	FromString string `json:"fromString"`
	ToString   string `json:"toString"`
}
