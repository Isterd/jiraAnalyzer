package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"jiraAnalyzer/backend/internal/models"
	"jiraAnalyzer/backend/internal/repository"
)

type AnalyticsService struct {
	repo *repository.Repository
}

func NewAnalyticsService(repo *repository.Repository) *AnalyticsService {
	return &AnalyticsService{
		repo: repo,
	}
}

// Метод для получения аналитики
func (s *AnalyticsService) GetProjectAnalytics(ctx context.Context, projectKey string) (models.ProjectAnalytics, error) {
	if projectKey == "" {
		return models.ProjectAnalytics{}, &models.InvalidInputError{Message: "project key cannot be empty"}
	}

	analytics, err := s.repo.GetProjectAnalytics(ctx, projectKey)
	if err != nil {
		return models.ProjectAnalytics{}, fmt.Errorf("failed to get project analytics: %w", err)
	}

	return analytics, nil
}

func (s *AnalyticsService) IsProjectAnalyzed(ctx context.Context, projectKey string) (bool, error) {
	return s.repo.IsProjectAnalyzed(ctx, projectKey)
}

func (s *AnalyticsService) DeleteProjectAnalytics(ctx context.Context, projectKey string) error {
	return s.repo.DeleteProjectAnalytics(ctx, projectKey)
}

func (s *AnalyticsService) CalculateOpenTimeHistogram(ctx context.Context, projectKey string) ([]models.HistogramData, error) {
	histogram, err := s.repo.CalculateOpenTimeHistogram(ctx, projectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate histogram: %w", err)
	}

	data, err := json.Marshal(histogram)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal histogram data: %w", err)
	}

	if err := s.repo.SaveAnalytics(ctx, projectKey, 1, data); err != nil {
		return nil, fmt.Errorf("failed to save histogram data: %w", err)
	}

	return histogram, nil
}

func (s *AnalyticsService) GetOpenTimeHistogram(ctx context.Context, projectKey string) ([]models.HistogramData, error) {
	data, err := s.repo.GetAnalytics(ctx, projectKey, 1)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to get histogram data: %w", err)
	}

	if len(data) == 0 {
		return s.CalculateOpenTimeHistogram(ctx, projectKey)
	}

	var histogram []models.HistogramData
	if err := json.Unmarshal(data, &histogram); err != nil {
		return nil, fmt.Errorf("failed to unmarshal histogram data: %w", err)
	}

	return histogram, nil
}

func (s *AnalyticsService) CalculateStatusTimeDistribution(ctx context.Context, projectKey string) ([]models.StatusTimeData, error) {
	distribution, err := s.repo.CalculateStatusTimeDistribution(ctx, projectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate status time distribution: %w", err)
	}

	data, err := json.Marshal(distribution)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal status time distribution data: %w", err)
	}

	if err := s.repo.SaveAnalytics(ctx, projectKey, 2, data); err != nil {
		return nil, fmt.Errorf("failed to save status time distribution data: %w", err)
	}

	return distribution, nil
}

func (s *AnalyticsService) GetStatusTimeDistribution(ctx context.Context, projectKey string) ([]models.StatusTimeData, error) {
	data, err := s.repo.GetAnalytics(ctx, projectKey, 2)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to get status time distribution data: %w", err)
	}

	if len(data) == 0 {
		return s.CalculateStatusTimeDistribution(ctx, projectKey)
	}

	var distribution []models.StatusTimeData
	if err := json.Unmarshal(data, &distribution); err != nil {
		return nil, fmt.Errorf("failed to unmarshal status time distribution data: %w", err)
	}

	return distribution, nil
}

func (s *AnalyticsService) CalculateActivityGraph(ctx context.Context, projectKey string) ([]models.ActivityData, error) {
	activity, err := s.repo.CalculateActivityGraph(ctx, projectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate activity graph: %w", err)
	}

	data, err := json.Marshal(activity)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal activity graph data: %w", err)
	}

	if err := s.repo.SaveAnalytics(ctx, projectKey, 3, data); err != nil {
		return nil, fmt.Errorf("failed to save activity graph data: %w", err)
	}

	return activity, nil
}

func (s *AnalyticsService) GetActivityGraph(ctx context.Context, projectKey string) ([]models.ActivityData, error) {
	data, err := s.repo.GetAnalytics(ctx, projectKey, 3)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to get activity graph data: %w", err)
	}

	if len(data) == 0 {
		return s.CalculateActivityGraph(ctx, projectKey)
	}

	var activity []models.ActivityData
	if err := json.Unmarshal(data, &activity); err != nil {
		return nil, fmt.Errorf("failed to unmarshal activity graph data: %w", err)
	}

	return activity, nil
}

func (s *AnalyticsService) CalculateComplexityGraph(ctx context.Context, projectKey string) ([]models.ComplexityData, error) {
	complexity, err := s.repo.CalculateComplexityGraph(ctx, projectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate complexity graph: %w", err)
	}

	data, err := json.Marshal(complexity)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal complexity graph data: %w", err)
	}

	if err := s.repo.SaveAnalytics(ctx, projectKey, 5, data); err != nil {
		return nil, fmt.Errorf("failed to save complexity graph data: %w", err)
	}

	return complexity, nil
}

func (s *AnalyticsService) GetComplexityGraph(ctx context.Context, projectKey string) ([]models.ComplexityData, error) {
	data, err := s.repo.GetAnalytics(ctx, projectKey, 5)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to get complexity graph data: %w", err)
	}

	if len(data) == 0 {
		return s.CalculateComplexityGraph(ctx, projectKey)
	}

	var complexity []models.ComplexityData
	if err := json.Unmarshal(data, &complexity); err != nil {
		return nil, fmt.Errorf("failed to unmarshal complexity graph data: %w", err)
	}

	return complexity, nil
}

func (s *AnalyticsService) CalculatePriorityDistribution(ctx context.Context, projectKey string) ([]models.PriorityData, error) {
	distribution, err := s.repo.CalculatePriorityDistribution(ctx, projectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate priority distribution: %w", err)
	}

	data, err := json.Marshal(distribution)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal priority distribution data: %w", err)
	}

	if err := s.repo.SaveAnalytics(ctx, projectKey, 5, data); err != nil {
		return nil, fmt.Errorf("failed to save priority distribution data: %w", err)
	}

	return distribution, nil
}

func (s *AnalyticsService) GetPriorityDistribution(ctx context.Context, projectKey string) ([]models.PriorityData, error) {
	data, err := s.repo.GetAnalytics(ctx, projectKey, 5)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to get priority distribution data: %w", err)
	}

	if len(data) == 0 {
		return s.CalculatePriorityDistribution(ctx, projectKey)
	}

	var distribution []models.PriorityData
	if err := json.Unmarshal(data, &distribution); err != nil {
		return nil, fmt.Errorf("failed to unmarshal priority distribution data: %w", err)
	}

	return distribution, nil
}

func (s *AnalyticsService) CalculatePriorityDistributionClosedTasks(ctx context.Context, projectKey string) ([]models.PriorityData, error) {
	distribution, err := s.repo.CalculatePriorityDistributionClosedTasks(ctx, projectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate priority distribution closed tasks: %w", err)
	}

	data, err := json.Marshal(distribution)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal priority distribution closed tasks data: %w", err)
	}

	if err := s.repo.SaveAnalytics(ctx, projectKey, 5, data); err != nil {
		return nil, fmt.Errorf("failed to save priority distribution closed tasks data: %w", err)
	}

	return distribution, nil
}

func (s *AnalyticsService) GetPriorityDistributionClosedTasks(ctx context.Context, projectKey string) ([]models.PriorityData, error) {
	data, err := s.repo.GetAnalytics(ctx, projectKey, 5)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to get priority distribution closed tasks data: %w", err)
	}

	if len(data) == 0 {
		return s.CalculatePriorityDistributionClosedTasks(ctx, projectKey)
	}

	var distribution []models.PriorityData
	if err := json.Unmarshal(data, &distribution); err != nil {
		return nil, fmt.Errorf("failed to unmarshal priority distribution closed tasks data: %w", err)
	}

	return distribution, nil
}

func (s *AnalyticsService) CalculateComparison(ctx context.Context, projectKey1, projectKey2 string) ([]models.ComparisonProjects, error) {
	comparison, err := s.repo.CompareProjects(ctx, projectKey1, projectKey2)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate project comparison: %w", err)
	}

	data, err := json.Marshal(comparison)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal comparison data: %w", err)
	}

	if err := s.repo.SaveAnalytics(ctx, fmt.Sprintf("%s-%s", projectKey1, projectKey2), 6, data); err != nil {
		return nil, fmt.Errorf("failed to save comparison data: %w", err)
	}

	return comparison, nil
}

func (s *AnalyticsService) GetComparison(ctx context.Context, projectKey1, projectKey2 string) ([]models.ComparisonProjects, error) {
	data, err := s.repo.GetAnalytics(ctx, fmt.Sprintf("%s-%s", projectKey1, projectKey2), 6)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to get comparison data: %w", err)
	}

	if len(data) == 0 {
		return s.CalculateComparison(ctx, projectKey1, projectKey2)
	}

	var comparison []models.ComparisonProjects
	if err := json.Unmarshal(data, &comparison); err != nil {
		return nil, fmt.Errorf("failed to unmarshal comparison data: %w", err)
	}

	return comparison, nil
}
