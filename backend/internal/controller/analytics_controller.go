package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"jiraAnalyzer/backend/internal/config"
	"jiraAnalyzer/backend/internal/models"
	"jiraAnalyzer/backend/internal/service"
	"jiraAnalyzer/backend/internal/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AnalyticsController struct {
	service *service.AnalyticsService
	cfg     *config.Backend
}

func NewAnalyticsController(service *service.AnalyticsService, cfg config.Backend) *AnalyticsController {
	return &AnalyticsController{service: service, cfg: &cfg}
}

func (h *AnalyticsController) GetProjectAnalytics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectKey := vars["key"]

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	analytics, err := h.service.GetProjectAnalytics(ctx, projectKey)
	if err != nil {
		log.Printf("Error fetching analytics for project %s: %v", projectKey, err)

		// Проверяем тип ошибки и возвращаем соответствующий HTTP-статус
		var notFoundErr *models.NotFoundError // Предполагается, что у вас есть такой тип ошибки
		if errors.As(err, &notFoundErr) {
			http.Error(w, "Project not found", http.StatusNotFound)
			return
		}

		var invalidInputErr *models.InvalidInputError // Предполагается, что у вас есть такой тип ошибки
		if errors.As(err, &invalidInputErr) {
			http.Error(w, "Invalid project key", http.StatusBadRequest)
			return
		}

		// Если ошибка неизвестного типа, возвращаем 500 с подробным сообщением
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("internal server error: %w", err))
		return
	}

	response := map[string]interface{}{
		"data": analytics,
	}

	utils.WriteJSONResponse(w, response)
}

func (h *AnalyticsController) IsProjectAnalytics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectKey := vars["key"]

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	isAnalyzed, err := h.service.IsProjectAnalyzed(ctx, projectKey)
	if err != nil {
		log.Printf("Error fetching analytics for project %s: %v", projectKey, err)

		// Проверяем тип ошибки и возвращаем соответствующий HTTP-статус
		var notFoundErr *models.NotFoundError // Предполагается, что у вас есть такой тип ошибки
		if errors.As(err, &notFoundErr) {
			http.Error(w, "Project not found", http.StatusNotFound)
			return
		}

		var invalidInputErr *models.InvalidInputError // Предполагается, что у вас есть такой тип ошибки
		if errors.As(err, &invalidInputErr) {
			http.Error(w, "Invalid project key", http.StatusBadRequest)
			return
		}

		// Если ошибка неизвестного типа, возвращаем 500 с подробным сообщением
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("internal server error: %w", err))
		return
	}

	response := map[string]interface{}{
		"is_analyzed": isAnalyzed,
	}

	utils.WriteJSONResponse(w, response)
}

func (h *AnalyticsController) DeleteProjectAnalytics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectKey := vars["key"]

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.service.DeleteProjectAnalytics(ctx, projectKey); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("error while deleting analytics: %w", err))
		return
	}

	response := map[string]interface{}{
		"success": fmt.Sprintf("all analytics data for project %s has been deleted", projectKey),
	}

	utils.WriteJSONResponse(w, response)
}

// GET /api/v1/graph/get/{taskNumber:[0-9]+}
func (h *AnalyticsController) GetGraph(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	taskNumber, err := strconv.Atoi(params["taskNumber"])
	if err != nil || taskNumber < 1 || taskNumber > 6 {
		http.Error(w, "Invalid task number", http.StatusBadRequest)
		return
	}

	projectKey := r.URL.Query().Get("project")
	if projectKey == "" {
		http.Error(w, "Missing project key", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.AnalyticsTimeout)
	defer cancel()

	data, err := h.getAnalyticsData(ctx, taskNumber, projectKey)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "Request Timeout", http.StatusRequestTimeout)
			return
		}
		http.Error(w, fmt.Sprintf("Error processing project %s: %v", projectKey, err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"_links": map[string]string{
			"self": fmt.Sprintf("/api/v1/graph/get/%d?project=%s", taskNumber, projectKey),
		},
		"data": data,
	}

	utils.WriteJSONResponse(w, response)
}

// POST /api/v1/graph/make/{taskNumber:[0-9]}
func (h *AnalyticsController) MakeGraph(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	taskNumber, err := strconv.Atoi(params["taskNumber"])
	if err != nil || taskNumber < 1 || taskNumber > 6 {
		http.Error(w, "invalid task number", http.StatusBadRequest)
		return
	}

	projectKey := r.URL.Query().Get("project")
	if projectKey == "" {
		http.Error(w, "missing project key", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.AnalyticsTimeout)
	defer cancel()

	switch taskNumber {
	case 1:
		_, err = h.service.CalculateOpenTimeHistogram(ctx, projectKey, taskNumber)
	case 2:
		_, err = h.service.CalculateStatusTimeDistribution(ctx, projectKey, taskNumber)
	case 3:
		_, err = h.service.CalculateActivityGraph(ctx, projectKey, taskNumber)
	case 4:
		_, err = h.service.CalculateComplexityGraph(ctx, projectKey, taskNumber)
	case 5:
		_, err = h.service.CalculatePriorityDistribution(ctx, projectKey, taskNumber)
	case 6:
		_, err = h.service.CalculatePriorityDistributionClosedTasks(ctx, projectKey, taskNumber)
	default:
		http.Error(w, "invalid task number", http.StatusBadRequest)
		return
	}

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			http.Error(w, "Request Timeout", http.StatusRequestTimeout)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"_links": map[string]string{
			"self": fmt.Sprintf("/api/v1/graph/make/%d?project=%s", taskNumber, projectKey),
		},
		"success": "Graph calculation started",
	}

	utils.WriteJSONResponse(w, response)
}

// GET /api/v1/compare/{taskNumber:[0-9]+}
func (h *AnalyticsController) GetComparison(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	taskNumber, err := strconv.Atoi(params["taskNumber"])
	if err != nil || taskNumber < 1 || taskNumber > 6 {
		http.Error(w, "Invalid task number", http.StatusBadRequest)
		return
	}

	projectKeys := strings.Split(r.URL.Query().Get("project"), ",")
	if len(projectKeys) < 2 || len(projectKeys) > 3 {
		http.Error(w, "Between 2 and 3 project keys are required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.AnalyticsTimeout)
	defer cancel()

	var comparisonResults []map[string]interface{}

	for _, projectKey := range projectKeys {
		data, err := h.getAnalyticsData(ctx, taskNumber, projectKey)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				http.Error(w, "Request Timeout", http.StatusRequestTimeout)
				return
			}
			http.Error(w, fmt.Sprintf("Error processing project %s: %v", projectKey, err), http.StatusInternalServerError)
			return
		}

		comparisonResults = append(comparisonResults, map[string]interface{}{
			"project_key": projectKey,
			"data":        data,
		})
	}

	response := map[string]interface{}{
		"_links": map[string]string{
			"self": fmt.Sprintf("/api/v1/compare/%d?project=%s", taskNumber, strings.Join(projectKeys, ",")),
		},
		"data": comparisonResults,
	}

	utils.WriteJSONResponse(w, response)
}

func (h *AnalyticsController) getAnalyticsData(ctx context.Context, taskNumber int, projectKey string) (interface{}, error) {
	switch taskNumber {
	case 1:
		return h.service.GetOpenTimeHistogram(ctx, projectKey, taskNumber)
	case 2:
		return h.service.GetStatusTimeDistribution(ctx, projectKey, taskNumber)
	case 3:
		return h.service.GetActivityGraph(ctx, projectKey, taskNumber)
	case 4:
		return h.service.GetComplexityGraph(ctx, projectKey, taskNumber)
	case 5:
		return h.service.GetPriorityDistribution(ctx, projectKey, taskNumber)
	case 6:
		return h.service.GetPriorityDistributionClosedTasks(ctx, projectKey, taskNumber)
	default:
		return nil, fmt.Errorf("invalid task number")
	}
}
