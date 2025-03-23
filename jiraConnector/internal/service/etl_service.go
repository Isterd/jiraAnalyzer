package service

import (
	"context"
	"fmt"
	"jiraAnalyzer/jiraConnector/internal/models"
	"jiraAnalyzer/jiraConnector/internal/repository/database"
	"log"
	"strings"
	"sync"
	"time"
)

type ETLService struct {
	jiraService       *JiraService
	repo              *database.Repository
	ThreadCount       int
	IssueInOneRequest int
}

func NewETLService(jiraService *JiraService, repo *database.Repository, threadCount int, issueInOneRequest int) *ETLService {
	return &ETLService{
		jiraService:       jiraService,
		repo:              repo,
		ThreadCount:       threadCount,
		IssueInOneRequest: issueInOneRequest,
	}
}

func (s *ETLService) GetProjectsFromJira(ctx context.Context, page, limit int, search string) ([]models.DBProject, models.PageInfo, error) {
	// Получаем все проекты из Jira
	jiraProjects, err := s.jiraService.GetAllProjects(ctx)
	if err != nil {
		return nil, models.PageInfo{}, fmt.Errorf("failed to fetch projects from Jira: %w", err)
	}

	// Фильтруем проекты по параметру search
	var filteredProjects []models.DBProject
	for _, project := range jiraProjects {
		if search == "" || strings.Contains(strings.ToLower(project.Name), strings.ToLower(search)) || strings.Contains(strings.ToLower(project.Key), strings.ToLower(search)) {
			filteredProjects = append(filteredProjects, models.DBProject{
				Key:  project.Key,
				Name: project.Name,
				URL:  project.URL,
			})
		}
	}

	// Вычисляем общее количество проектов
	totalCount := len(filteredProjects)

	// Пагинация
	offset := (page - 1) * limit
	end := offset + limit
	if end > totalCount {
		end = totalCount
	}
	paginatedProjects := filteredProjects[offset:end]

	// Формируем информацию о пагинации
	pageInfo := models.PageInfo{
		CurrentPage: page,
		PageCount:   (totalCount + limit - 1) / limit,
		TotalCount:  totalCount,
	}

	return paginatedProjects, pageInfo, nil
}

func (s *ETLService) UpdateProject(ctx context.Context, projectKeys []string) error {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errChan := make(chan error, 1)

	for _, projectKey := range projectKeys {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			if err := s.updateSingleProject(ctx, key); err != nil {
				select {
				case errChan <- fmt.Errorf("failed to update project %s: %w", key, err):
				default:
				}
				cancel() // Прерываем все горутины при ошибке
			}
		}(projectKey)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *ETLService) updateSingleProject(ctx context.Context, projectKey string) error {
	exists, err := s.repo.CheckProjectExists(ctx, projectKey)
	if err != nil {
		return fmt.Errorf("failed to check project existence: %w", err)
	}

	// Если проект новый - загружаем метаданные
	if !exists {
		log.Printf("Project %s not found in DB, fetching metadata...", projectKey)
		projects, err := s.jiraService.GetAllProjects(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch projects: %w", err)
		}
		project, err := s.transformProject(ctx, projects, projectKey)
		if err != nil {
			return fmt.Errorf("failed to transform project: %w", err)
		}
		if err := s.repo.SaveProject(ctx, project); err != nil {
			return fmt.Errorf("failed to save project: %w", err)
		}
	}

	// Загружаем issues с адаптивной обработкой рейт-лимитов
	log.Printf("Loading issues for project %s...", projectKey)
	return s.loadIssuesWithBackoff(ctx, projectKey)
}

func (s *ETLService) loadIssuesWithBackoff(ctx context.Context, projectKey string) error {
	sem := make(chan struct{}, s.ThreadCount) // Semaphore to limit goroutines
	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	totalIssues := s.getIssueCount(projectKey)
	if totalIssues == 0 {
		log.Printf("No issues found for project %s", projectKey)
		return nil
	}
	batches := (totalIssues + s.IssueInOneRequest - 1) / s.IssueInOneRequest

	for i := 0; i < batches; i++ {
		if ctx.Err() != nil {
			break
		}

		wg.Add(1)
		sem <- struct{}{} // Acquire semaphore
		go func(startAt int) {
			defer func() {
				wg.Done()
				<-sem
			}()

			if err := s.loadIssuesBatch(ctx, projectKey, startAt); err != nil {
				select {
				case errChan <- fmt.Errorf("failed to load batch: %w", err):
				default:
				}

				cancel() // Прерываем все горутины при ошибке
			}
		}(i * s.IssueInOneRequest)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *ETLService) loadIssuesBatch(ctx context.Context, projectKey string, startAt int) error {
	log.Printf("Loading batch for project %s starting at %d", projectKey, startAt)

	issues, err := s.jiraService.GetProjectIssues(ctx, projectKey, startAt)
	if err != nil {
		return fmt.Errorf("failed to get project issues: %w", err)
	}
	log.Printf("Fetched %d issues for project %s starting at %d", len(issues), projectKey, startAt)

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	dbIssues := make([]models.DBIssue, len(issues))
	dbChangelogs := make([]models.DBChangelog, 0)

	for i, issue := range issues {
		log.Printf("Transforming issue: %s", issue.Key)
		dbIssues[i], err = s.transformIssue(issue, projectKey)
		if err != nil {
			return fmt.Errorf("failed to transform issue: %w", err)
		}
		changelogs, err := s.extractChangelogs(issue)
		if err != nil {
			return fmt.Errorf("failed to extract changelogs: %w", err)
		}
		dbChangelogs = append(dbChangelogs, changelogs...)
	}

	if err := s.repo.SaveIssuesTx(tx, dbIssues); err != nil {
		return fmt.Errorf("failed to save issues: %w", err)
	}

	if err := s.repo.SaveChangelogTx(tx, dbChangelogs); err != nil {
		return fmt.Errorf("failed to save changelogs: %w", err)
	}

	return tx.Commit()
}

func (s *ETLService) getIssueCount(projectKey string) int {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	count, err := s.jiraService.GetIssueCount(ctx, projectKey)
	if err != nil {
		log.Printf("Failed to get issue count for project %s: %v", projectKey, err)
		return 0
	}

	return count
}
