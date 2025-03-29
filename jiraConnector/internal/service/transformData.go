package service

import (
	"fmt"
	"jiraAnalyzer/jiraConnector/internal/models"
	"log"
	"strings"
	"time"
)

func (s *ETLService) transformProject(projects []models.JiraProject, projectKey string) (models.DBProject, error) {
	for _, p := range projects {
		if p.Key == projectKey {
			return models.DBProject{
				Key:  p.Key,
				Name: p.Name,
				URL:  p.URL,
			}, nil
		}
	}
	return models.DBProject{}, fmt.Errorf("project %s not found", projectKey)
}

func (s *ETLService) transformIssue(issue models.JiraIssue, projectKey string) (models.DBIssue, error) {
	dbIssue := models.DBIssue{}

	creatorID, err := s.repo.GetOrCreateAuthor(issue.Fields.Creator.DisplayName)
	if err != nil {
		return dbIssue, fmt.Errorf("failed to get or create creator: %w", err)
	}

	var assigneeID *int
	if issue.Fields.Assignee != nil {
		log.Printf("Assignee: %+v", issue.Fields.Assignee)
		id, err := s.repo.GetOrCreateAuthor(issue.Fields.Assignee.DisplayName)
		if err != nil {
			return dbIssue, fmt.Errorf("failed to get or create assignee: %w", err)
		}
		assigneeID = &id
	}

	var closedTime *time.Time
	if parsedTime, err := GetClosedTime(issue.Changelog); err == nil {
		closedTime = parsedTime
	} else {
		log.Printf("Failed to get closed time for issue %s: %v", issue.Key, err)
	}

	log.Printf("Transforming issue: %s, creator: %v, assignee: %v, timespent: %d, closed: %v",
		issue.Key,
		issue.Fields.Creator.DisplayName,
		issue.Fields.Assignee,
		issue.Fields.TimeSpent,
		func() string {
			if closedTime != nil {
				return closedTime.Format(time.RFC3339)
			}
			return "NULL"
		}(),
	)

	dbIssue = models.DBIssue{
		Key:         issue.Key,
		ProjectKey:  projectKey,
		Created:     parseJiraTime(issue.Fields.Created),
		Updated:     parseJiraTime(issue.Fields.Updated),
		Closed:      closedTime,
		Summary:     issue.Fields.Summary,
		Description: issue.Fields.Description,
		Type:        issue.Fields.IssueType.Name,
		Priority:    issue.Fields.Priority.Name,
		Status:      issue.Fields.Status.Name,
		TimeSpent:   issue.Fields.TimeSpent,
		CreatorID:   creatorID,
		AssigneeID:  assigneeID,
	}

	return dbIssue, nil
}

func (s *ETLService) extractChangelogs(issue models.JiraIssue) ([]models.DBChangelog, error) {
	var dbChangelogs []models.DBChangelog
	for _, history := range issue.Changelog.Histories {
		for _, item := range history.Items {
			if item.Field == "status" {
				authorID, err := s.repo.GetOrCreateAuthor(history.Author.DisplayName)
				if err != nil {
					return nil, err
				}

				log.Printf("Processing changelog: from=%s, to=%s", item.FromString, item.ToString)
				dbChangelogs = append(dbChangelogs, models.DBChangelog{
					IssueID:    issue.Key,
					AuthorID:   authorID,
					Created:    parseJiraTime(history.Created),
					FromStatus: item.FromString,
					ToStatus:   item.ToString,
				})
			}
		}
	}

	return dbChangelogs, nil
}

func parseJiraTime(str string) time.Time {
	layout := "2006-01-02T15:04:05.000-0700"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Printf("Failed to parse Jira time: %v", err)
		return time.Time{} // Возвращаем пустую дату
	}
	return t
}

func GetClosedTime(changelog models.JiraChangelog) (*time.Time, error) {
	log.Printf("Processing changelog with %d histories", len(changelog.Histories))
	for i := len(changelog.Histories) - 1; i >= 0; i-- {
		history := changelog.Histories[i]
		log.Printf("History created at: %s", history.Created)
		for _, item := range history.Items {
			log.Printf("Processing changelog: from=%s, to=%s", item.FromString, item.ToString)
			if item.Field == "status" && strings.ToLower(item.ToString) == "closed" {
				parsedTime := parseJiraTime(history.Created)
				if parsedTime.IsZero() {
					return nil, fmt.Errorf("failed to parse closed time")
				}
				return &parsedTime, nil
			}
		}
	}
	return nil, fmt.Errorf("task is not closed")
}
