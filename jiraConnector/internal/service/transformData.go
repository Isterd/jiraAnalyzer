package service

import (
	"context"
	"fmt"
	"jiraAnalyzer/jiraConnector/internal/models"
	"time"
)

func (s *ETLService) transformProject(ctx context.Context, projects []models.JiraProject, projectKey string) (models.DBProject, error) {
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
		return dbIssue, err
	}

	var assigneeID *int
	if issue.Fields.Assignee != nil {
		id, err := s.repo.GetOrCreateAuthor(issue.Fields.Assignee.DisplayName)
		if err != nil {
			return dbIssue, err
		}
		assigneeID = &id
	}

	dbIssue = models.DBIssue{
		JiraID:         issue.Key,
		ProjectKey:     projectKey,
		Key:            issue.Key,
		Created:        parseJiraTime(issue.Fields.Created),
		Updated:        parseJiraTime(issue.Fields.Updated),
		ResolutionDate: parseResolutionDate(issue.Fields.Resolution),
		Summary:        issue.Fields.Summary,
		Description:    issue.Fields.Description,
		Type:           issue.Fields.IssueType.Name,
		Priority:       issue.Fields.Priority.Name,
		Status:         issue.Fields.Status.Name,
		TimeSpent:      issue.Fields.TimeSpent,
		CreatorID:      creatorID,
		AssigneeID:     assigneeID,
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
	t, _ := time.Parse("2006-01-02T15:04:05Z", str)
	return t
}

func parseResolutionDate(res *models.JiraResolution) *time.Time {
	if res == nil || res.Date == "" {
		return nil
	}
	t := parseJiraTime(res.Date)
	return &t
}
