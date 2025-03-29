package database

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"jiraAnalyzer/backend/internal/models"
	"log"
)

type AnalyticsPostgres struct {
	db *sqlx.DB
}

func NewAnalyticsPostgres(db *sqlx.DB) *AnalyticsPostgres {
	return &AnalyticsPostgres{db: db}
}

func (r *AnalyticsPostgres) GetProjectAnalytics(ctx context.Context, projectKey string) (models.ProjectAnalytics, error) {
	var analytics models.ProjectAnalytics

	// Общее количество задач
	err := r.db.GetContext(ctx, &analytics.TotalIssues, `
        SELECT COUNT(*) FROM issues WHERE project_key = $1
    `, projectKey)
	if err != nil {
		return analytics, fmt.Errorf("failed to get total issues: %w", err)
	}

	// Количество закрытых задач
	err = r.db.GetContext(ctx, &analytics.ClosedIssues, `
        SELECT COUNT(*) FROM issues WHERE project_key = $1 AND status = 'Closed'
    `, projectKey)
	if err != nil {
		return analytics, fmt.Errorf("failed to get closed issues: %w", err)
	}

	// Шаг 3: Получение количества открытых задач
	err = r.db.GetContext(ctx, &analytics.OpenIssues, `
        SELECT COUNT(*) 
        FROM issues 
        WHERE project_key = $1 AND status = 'Open'
    `, projectKey)
	if err != nil {
		return analytics, fmt.Errorf("failed to fetch open issues count: %w", err)
	}

	// Количество переоткрытых задач
	err = r.db.GetContext(ctx, &analytics.ReopenIssues, `
		SELECT COUNT(DISTINCT issue_id) AS reopened_tasks 
		FROM status_changes WHERE issue_id IN (
    	SELECT key FROM issues WHERE project_key = $1
    	) AND from_status IN ('Closed', 'Resolved')
		AND to_status NOT IN ('Closed', 'Resolved')
		`, projectKey)
	if err != nil {
		return analytics, fmt.Errorf("failed to get reopen issues count: %w", err)
	}

	// Количество разрешенных задач
	err = r.db.GetContext(ctx, &analytics.ResolvedIssues, `SELECT COUNT(*) AS resolved_tasks
		FROM issues
		WHERE project_key = $1 AND status = 'Resolved'
		`, projectKey)
	if err != nil {
		return analytics, fmt.Errorf("failed to get resolved issues count: %w", err)
	}

	err = r.db.GetContext(ctx, &analytics.InProgressIssues, `SELECT COUNT(*) AS in_progress_tasks
		FROM issues
		WHERE project_key = $1 AND status = 'In Progress'
		`, projectKey)
	if err != nil {
		return analytics, fmt.Errorf("failed to get in_progress issues count: %w", err)
	}

	// Среднее время выполнения задачи (часы)
	err = r.db.GetContext(ctx, &analytics.AverageTimeIssues, `
		SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (closed - created)) / 3600), 0)
    	AS avg_completion_time_hours
		FROM issues
		WHERE project_key = $1 AND closed IS NOT NULL`, projectKey)
	if err != nil {
		return analytics, fmt.Errorf("failed to get average time issues: %w", err)
	}

	// Среднее количество заведенных задач в день за последнюю неделю
	err = r.db.GetContext(ctx, &analytics.AverageCountIssues, `
		SELECT COUNT(*) / 7 
    	AS avg_tasks_per_day_last_week
		FROM issues
		WHERE project_key = $1
		AND created >= NOW() - INTERVAL '7 days'
		`, projectKey)
	if err != nil {
		return analytics, fmt.Errorf("failed to get average issues count in last week: %w", err)
	}

	log.Printf("Analytics calculation completed successfully for project: %s", projectKey)
	return analytics, nil
}

func (r *AnalyticsPostgres) SaveAnalytics(ctx context.Context, projectKey string, taskNumber int, data []byte) error {
	query := `
        INSERT INTO analytics (project_key, task_number, data)
        VALUES ($1, $2, $3)
        ON CONFLICT (project_key, task_number) DO UPDATE
        SET data = EXCLUDED.data, created_at = NOW()
    `
	_, err := r.db.ExecContext(ctx, query, projectKey, taskNumber, data)
	return err
}

func (r *AnalyticsPostgres) GetAnalytics(ctx context.Context, projectKey string, taskNumber int) ([]byte, error) {
	var data []byte
	query := "SELECT data FROM analytics WHERE project_key = $1 AND task_number = $2 ORDER BY created_at DESC LIMIT 1"
	err := r.db.QueryRowContext(ctx, query, projectKey, taskNumber).Scan(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to get analytics: %w", err)
	}
	return data, nil
}

func (r *AnalyticsPostgres) IsProjectAnalyzed(ctx context.Context, projectKey string) (bool, error) {
	query := `
        SELECT EXISTS(
            SELECT 1 FROM analytics WHERE project_key = $1
        )
    `
	var exists bool
	err := r.db.QueryRowContext(ctx, query, projectKey).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if project is analyzed: %w", err)
	}
	return exists, nil
}

func (r *AnalyticsPostgres) DeleteProjectAnalytics(ctx context.Context, projectKey string) error {
	query := `
        DELETE FROM analytics WHERE project_key = $1
    `
	_, err := r.db.ExecContext(ctx, query, projectKey)
	if err != nil {
		return fmt.Errorf("failed to delete analytics for project: %w", err)
	}
	return nil
}

func (r *AnalyticsPostgres) CalculateOpenTimeHistogram(ctx context.Context, projectKey string) ([]models.HistogramData, error) {
	query := `
    WITH task_durations AS (
    	SELECT EXTRACT(EPOCH FROM (closed - created)) / 3600 AS duration_hours
    	FROM issues
    	WHERE project_key = $1 AND closed IS NOT NULL
	)
	SELECT 
    	FLOOR(duration_hours / 24) AS day_interval,
    	COUNT(*) AS task_count
	FROM task_durations
	GROUP BY day_interval
	ORDER BY day_interval
    `

	var histogram []models.HistogramData
	err := r.db.SelectContext(ctx, &histogram, query, projectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate open time histogram: %w", err)
	}
	return histogram, nil
}

func (r *AnalyticsPostgres) CalculateStatusTimeDistribution(ctx context.Context, projectKey string) ([]models.StatusTimeData, error) {
	query := `
    WITH status_durations AS (
    	SELECT 
        	issue_id,
        	from_status,
        	to_status,
        	EXTRACT(EPOCH FROM (created - LAG(created) OVER (PARTITION BY issue_id ORDER BY created))) / 3600 AS duration_hours
    	FROM status_changes
    	WHERE issue_id IN (SELECT key FROM issues WHERE project_key = $1)
	)
	SELECT 
    	from_status AS status,
    	FLOOR(duration_hours / 24) AS day_interval,
    	COUNT(*) AS task_count
	FROM status_durations
	WHERE duration_hours IS NOT NULL
	GROUP BY from_status, day_interval
	ORDER BY from_status, day_interval
    `

	var distribution []models.StatusTimeData
	err := r.db.SelectContext(ctx, &distribution, query, projectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate status time distribution: %w", err)
	}
	return distribution, nil
}

func (r *AnalyticsPostgres) CalculateActivityGraph(ctx context.Context, projectKey string) ([]models.ActivityData, error) {
	query := `
    WITH daily_stats AS (
    	SELECT 
        	DATE_TRUNC('day', created) AS day,
        	COUNT(*) FILTER (WHERE status = 'Open') AS opened_tasks,
        	COUNT(*) FILTER (WHERE status = 'Closed') AS closed_tasks
    	FROM issues
    	WHERE project_key = $1
    	GROUP BY day
	)
	SELECT 
	    day,
    	SUM(opened_tasks) OVER (ORDER BY day) AS cumulative_opened,
    	SUM(closed_tasks) OVER (ORDER BY day) AS cumulative_closed
	FROM daily_stats
	ORDER BY day
    `

	var activity []models.ActivityData
	err := r.db.SelectContext(ctx, &activity, query, projectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate activity graph: %w", err)
	}
	return activity, nil
}

func (r *AnalyticsPostgres) CalculateComplexityGraph(ctx context.Context, projectKey string) ([]models.ComplexityData, error) {
	query := `
        SELECT FLOOR(EXTRACT(EPOCH FROM (closed - created)) / 3600) AS time_spent_hours, COUNT(*) AS task_count
        FROM issues
        WHERE project_key = $1 AND closed IS NOT NULL
        GROUP BY time_spent_hours
        ORDER BY time_spent_hours
    `
	var data []models.ComplexityData
	err := r.db.SelectContext(ctx, &data, query, projectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate complexity graph: %w", err)
	}

	return data, nil
}

func (r *AnalyticsPostgres) CalculatePriorityDistribution(ctx context.Context, projectKey string) ([]models.PriorityData, error) {
	query := `
        SELECT 
            priority, COUNT(*) AS task_count
        FROM issues
        WHERE project_key = $1
        GROUP BY priority
        ORDER BY task_count DESC;
    `

	var distribution []models.PriorityData
	err := r.db.SelectContext(ctx, &distribution, query, projectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate priority distribution: %w", err)
	}

	return distribution, nil
}

func (r *AnalyticsPostgres) CalculatePriorityDistributionClosedTasks(ctx context.Context, projectKey string) ([]models.PriorityData, error) {
	query := `
        SELECT priority, COUNT(*) AS task_count
        FROM issues
        WHERE project_key = $1 AND closed IS NOT NULL
        GROUP BY priority
        ORDER BY task_count DESC
    `
	var data []models.PriorityData
	err := r.db.SelectContext(ctx, &data, query, projectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate priority distribution for closed tasks: %w", err)
	}
	return data, nil
}

func (r *AnalyticsPostgres) CompareProjects(ctx context.Context, projectKey1, projectKey2 string) ([]models.ComparisonProjects, error) {
	query := `
	SELECT 
    	project_key,
    	COUNT(*) AS total_tasks,
    	COUNT(*) FILTER (WHERE status = 'Open') AS open_tasks,
    	COUNT(*) FILTER (WHERE status IN ('Closed', 'Resolved')) AS closed_tasks,
    	AVG(EXTRACT(EPOCH FROM (closed - created)) / 3600) AS avg_completion_time_hours
	FROM issues
	WHERE project_key IN ($1, $2)
	GROUP BY project_key
	`

	var projects []models.ComparisonProjects
	err := r.db.SelectContext(ctx, &projects, query, projectKey1, projectKey2)
	if err != nil {
		return nil, fmt.Errorf("failed to compare projects: %w", err)
	}

	return projects, nil
}
