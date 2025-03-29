package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"jiraAnalyzer/backend/internal/models"
	"net/http"
	"strings"
	"time"
)

type HTTPJiraClient struct {
	client *http.Client
	url    string
}

func NewHTTPJiraClient(url string) *HTTPJiraClient {
	return &HTTPJiraClient{
		client: &http.Client{Timeout: 30 * time.Second},
		url:    url,
	}
}

func (c *HTTPJiraClient) GetConnectorProjects(ctx context.Context, page, limit int, search string) ([]models.Project, models.PageInfo, error) {
	url := fmt.Sprintf("%s/projects?page=%d&limit=%d&search=%s", c.url, page, limit, search)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, models.PageInfo{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, models.PageInfo{}, fmt.Errorf("failed to call JiraConnector: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, models.PageInfo{}, fmt.Errorf("unexpected status code from JiraConnector: %d", resp.StatusCode)
	}

	var response struct {
		Projects []models.Project `json:"projects"`
		PageInfo models.PageInfo  `json:"pageInfo"`
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, models.PageInfo{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Projects, response.PageInfo, nil
}

func (c *HTTPJiraClient) UpdateConnectorProject(ctx context.Context, projectKeys []string) error {
	url := fmt.Sprintf("%s/updateProject?projects=%s", c.url, strings.Join(projectKeys, ","))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call JiraConnector: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from JiraConnector: %d", resp.StatusCode)
	}

	return nil
}
