package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"jiraAnalyzer/jiraConnector/internal/models"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type ClientConfig struct {
	JiraUrl           string        `yaml:"jiraUrl"`
	IssueInOneRequest int           `yaml:"issueInOneRequest"`
	ThreadCount       int           `yaml:"threadCount"`
	MaxAttempts       int           `yaml:"maxAttempts"`
	MaxTimeSleep      time.Duration `yaml:"maxTimeSleep"`
	MinTimeSleep      time.Duration `yaml:"minTimeSleep"`
}

type Jira struct {
	cfg        ClientConfig
	clientPool []*http.Client
}

func NewJiraClient(cfg ClientConfig) *Jira {
	clientPool := make([]*http.Client, cfg.ThreadCount)
	for i := 0; i < cfg.ThreadCount; i++ {
		clientPool[i] = &http.Client{Timeout: 60 * time.Second} //стоит убрать магическое число
	}

	return &Jira{
		cfg:        cfg,
		clientPool: clientPool,
	}
}

func (c *Jira) GetAllProjects(ctx context.Context) ([]models.JiraProject, error) {
	url := fmt.Sprintf("%s/rest/api/2/project", c.cfg.JiraUrl)
	var projects []models.JiraProject
	err := c.doRequestWithRetry(url, &projects, ctx)
	return projects, err
}

func (c *Jira) GetProjectIssues(ctx context.Context, projectKey string, startAt int) ([]models.JiraIssue, error) {
	url := fmt.Sprintf("%s/rest/api/2/search?jql=project=%s&startAt=%d&maxResults=%d&expand=changelog",
		c.cfg.JiraUrl, projectKey, startAt, c.cfg.IssueInOneRequest)

	var response models.JiraSearchResponse
	err := c.doRequestWithRetry(url, &response, ctx)
	if err != nil {
		log.Printf("Failed to fetch issues for project %s: %v", projectKey, err)
		return nil, fmt.Errorf("failed to fetch issues: %w", err)
	}
	return response.Issues, nil
}

func (c *Jira) GetIssueCount(ctx context.Context, projectKey string) (int, error) {
	url := fmt.Sprintf("%s/rest/api/2/search?jql=project=%s&maxResults=0", c.cfg.JiraUrl, projectKey)
	var response struct {
		Total int `json:"total"`
	}

	if err := c.doRequestWithRetry(url, &response, ctx); err != nil {
		return 0, fmt.Errorf("failed to get issue count for project %s: %w", projectKey, err)
	}

	return response.Total, nil
}

func (c *Jira) doRequestWithRetry(url string, response interface{}, ctx context.Context) error {
	attempt := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		clientIndex := rand.Intn(len(c.clientPool))
		req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
		resp, err := c.clientPool[clientIndex].Do(req)
		if err != nil {
			sleepTime := c.cfg.MinTimeSleep * time.Duration(math.Pow(2, float64(attempt)))
			if sleepTime > c.cfg.MaxTimeSleep || attempt >= c.cfg.MaxAttempts {
				return fmt.Errorf("max retries exceeded: %w", err)
			}
			time.Sleep(sleepTime)
			attempt++
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := resp.Header.Get("Retry-After")
			if retry, err := strconv.Atoi(retryAfter); err == nil {
				time.Sleep(time.Duration(retry) * time.Second)
				continue
			}

			sleepTime := c.cfg.MinTimeSleep * time.Duration(math.Pow(2, float64(attempt)))
			if sleepTime > c.cfg.MaxTimeSleep {
				return fmt.Errorf("rate limit exceeded: %v", err)
			}

			time.Sleep(sleepTime)
			attempt++
			continue
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("status: %d, body: %s", resp.StatusCode, body)
		}

		return json.NewDecoder(resp.Body).Decode(response)
	}
}
