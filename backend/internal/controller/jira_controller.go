package controller

import (
	"encoding/json"
	"jiraAnalyzer/backend/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type JiraController struct {
	service *service.JiraClientService
}

func NewJiraController(service *service.JiraClientService) *JiraController {
	return &JiraController{service: service}
}

func (h *JiraController) GetConnectorProjects(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	search := r.URL.Query().Get("search")

	projects, pageInfo, err := h.service.GetConnectorProjects(r.Context(), page, limit, search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"projects": projects,
		"pageInfo": pageInfo,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *JiraController) UpdateConnectorProject(w http.ResponseWriter, r *http.Request) {
	projectKeysParam := r.URL.Query().Get("projects")
	if projectKeysParam == "" {
		http.Error(w, "Missing project keys", http.StatusBadRequest)
		return
	}

	projectKeys := strings.Split(projectKeysParam, ",")
	err := h.service.UpdateConnectorProject(r.Context(), projectKeys)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success"}`))
}
