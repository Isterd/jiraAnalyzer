package jiraclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"jiraAnalyzer/jiraConnector/internal/backoff"
	"jiraAnalyzer/jiraConnector/internal/models"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type ProgramSettings struct {
	JiraUrl           string        `yaml:"jiraUrl"`
	IssueInOneRequest int           `yaml:"issueInOneRequest"`
	Host              string        `yaml:"bindAddress"`
	ThreadCount       int           `yaml:"threadCount"`
	MaxTimeSleep      time.Duration `yaml:"maxTimeSleep"`
	MinTimeSleep      time.Duration `yaml:"minTimeSleep"`
}

type JiraClient struct {
	cfg        ProgramSettings
	clientPool []*http.Client
	backOff    backoff.BackOff
}

func NewJiraClient(cfg ProgramSettings) *JiraClient {
	clientPool := make([]*http.Client, cfg.ThreadCount)
	for i := 0; i < cfg.ThreadCount; i++ {
		clientPool[i] = &http.Client{Timeout: 60 * time.Second}
	}

	return &JiraClient{
		cfg:        cfg,
		clientPool: clientPool,
		backOff:    backoff.NewExponentialBackOff(cfg.MinTimeSleep, cfg.MaxTimeSleep),
	}
}

func (c *JiraClient) GetAllProjects(ctx context.Context) ([]models.JiraProject, error) {
	url := fmt.Sprintf("%s/rest/api/2/project", c.cfg.JiraUrl)
	var projects []models.JiraProject
	err := c.doRequestWithRetry(url, &projects, ctx)
	return projects, err
}

func (c *JiraClient) GetProjectIssues(ctx context.Context, projectKey string, startAt int) ([]models.JiraIssue, error) {
	url := fmt.Sprintf("%s/rest/api/2/search?jql=project=%s&startAt=%d&maxResults=%d",
		c.cfg.JiraUrl, projectKey, startAt, c.cfg.IssueInOneRequest)

	var response models.JiraSearchResponse
	err := c.doRequestWithRetry(url, &response, ctx)
	if err != nil {
		log.Printf("Failed to fetch issues for project %s: %v", projectKey, err)
		return nil, fmt.Errorf("failed to fetch issues: %w", err)
	}
	return response.Issues, nil
}

func (c *JiraClient) GetIssueCount(ctx context.Context, projectKey string) (int, error) {
	url := fmt.Sprintf("%s/rest/api/2/search?jql=project=%s&maxResults=0", c.cfg.JiraUrl, projectKey)
	var response struct {
		Total int `json:"total"`
	}

	if err := c.doRequestWithRetry(url, &response, ctx); err != nil {
		return 0, fmt.Errorf("failed to get issue count for project %s: %w", projectKey, err)
	}

	return response.Total, nil
}

func (c *JiraClient) doRequestWithRetry(url string, response interface{}, ctx context.Context) error {
	var lastErr error
	sleepTime := c.cfg.MinTimeSleep

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Используем пул клиентов
		clientIndex := rand.Intn(len(c.clientPool))

		req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
		resp, err := c.clientPool[clientIndex].Do(req)
		if err != nil {
			lastErr = fmt.Errorf("network error: %v", err)
			time.Sleep(sleepTime)
			sleepTime *= 2
			if sleepTime > c.cfg.MaxTimeSleep {
				return lastErr
			}
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := resp.Header.Get("Retry-After")
			if retry, err := strconv.Atoi(retryAfter); err == nil {
				time.Sleep(time.Duration(retry) * time.Second)
				continue
			}
			time.Sleep(sleepTime)
			sleepTime *= 2
			if sleepTime > c.cfg.MaxTimeSleep {
				return fmt.Errorf("rate limit exceeded: %v", lastErr)
			}
			continue
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("status: %d, body: %s", resp.StatusCode, body)
		}

		return json.NewDecoder(resp.Body).Decode(response)
	}
}
