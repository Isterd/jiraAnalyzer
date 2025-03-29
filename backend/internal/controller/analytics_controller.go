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
		"message": fmt.Sprintf("all analytics data for project %s has been deleted", projectKey),
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

	var data interface{}

	switch taskNumber {
	case 1:
		data, err = h.service.GetOpenTimeHistogram(ctx, projectKey)
	case 2:
		data, err = h.service.GetStatusTimeDistribution(ctx, projectKey)
	case 3:
		data, err = h.service.GetActivityGraph(ctx, projectKey)
	case 4:
		data, err = h.service.GetComplexityGraph(ctx, projectKey)
	case 5:
		data, err = h.service.GetPriorityDistribution(ctx, projectKey)
	case 6:
		data, err = h.service.GetPriorityDistributionClosedTasks(ctx, projectKey)
	default:
		http.Error(w, "Invalid task number", http.StatusBadRequest)
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

	switch taskNumber {
	case 1:
		_, err = h.service.CalculateOpenTimeHistogram(ctx, projectKey)
	case 2:
		_, err = h.service.CalculateStatusTimeDistribution(ctx, projectKey)
	case 3:
		_, err = h.service.CalculateActivityGraph(ctx, projectKey)
	case 4:
		_, err = h.service.GetComplexityGraph(ctx, projectKey)
	case 5:
		_, err = h.service.GetPriorityDistribution(ctx, projectKey)
	case 6:
		_, err = h.service.GetPriorityDistributionClosedTasks(ctx, projectKey)
	default:
		http.Error(w, "Invalid task number", http.StatusBadRequest)
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
		"message": "Graph calculation started",
	}

	utils.WriteJSONResponse(w, response)
}

// GET /api/v1/compare/{taskNumber:[0-9]+}
func (h *AnalyticsController) GetComparison(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	taskNumber, err := strconv.Atoi(params["taskNumber"])
	if err != nil || taskNumber < 1 || taskNumber > 5 {
		http.Error(w, "Invalid task number", http.StatusBadRequest)
		return
	}

	projectKeys := strings.Split(r.URL.Query().Get("project"), ",")
	if len(projectKeys) != 2 {
		http.Error(w, "Exactly two project keys are required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.AnalyticsTimeout)
	defer cancel()

	var data interface{}

	switch taskNumber {
	case 5:
		data, err = h.service.GetComparison(ctx, projectKeys[0], projectKeys[1])
	default:
		http.Error(w, "Invalid task number", http.StatusBadRequest)
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
			"self": fmt.Sprintf("/api/v1/compare/%d?project=%s,%s", taskNumber, projectKeys[0], projectKeys[1]),
		},
		"data": data,
	}

	utils.WriteJSONResponse(w, response)
}

// POST /api/v1/compare/{taskNumber:[0-9]}
func (h *AnalyticsController) MakeComparison(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	taskNumber, err := strconv.Atoi(params["taskNumber"])
	if err != nil || taskNumber < 1 || taskNumber > 5 {
		http.Error(w, "Invalid task number", http.StatusBadRequest)
		return
	}

	projectKeys := strings.Split(r.URL.Query().Get("project"), ",")
	if len(projectKeys) != 2 {
		http.Error(w, "Exactly two project keys are required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.AnalyticsTimeout)
	defer cancel()

	switch taskNumber {
	case 5:
		_, err = h.service.CalculateComparison(ctx, projectKeys[0], projectKeys[1])
	default:
		http.Error(w, "Invalid task number", http.StatusBadRequest)
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
			"self": fmt.Sprintf("/api/v1/compare/%d?project=%s,%s", taskNumber, projectKeys[0], projectKeys[1]),
		},
		"message": "Project comparison calculation started",
	}

	utils.WriteJSONResponse(w, response)
}
