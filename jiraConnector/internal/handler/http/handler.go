package handler

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"jiraAnalyzer/jiraConnector/internal/models"
	"jiraAnalyzer/jiraConnector/internal/service"
	"strings"
	"time"

	"net/http"
	"strconv"
)

type JiraConnectorConfig struct {
	BaseUrl      string        `yaml:"baseUrl"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
}

type Handler struct {
	etlService *service.ETLService
	cfg        JiraConnectorConfig
}

func NewHandler(etlService *service.ETLService, r *mux.Router, cfg JiraConnectorConfig) *mux.Router {
	h := &Handler{etlService: etlService, cfg: cfg}

	r.HandleFunc("/updateProject", h.UpdateProject)
	r.HandleFunc("/projects", h.GetProjects).Methods(http.MethodOptions, http.MethodGet)

	return r
}

func (h *Handler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	projectKeysParam := r.URL.Query().Get("projects")
	if projectKeysParam == "" {
		http.Error(w, "Missing project keys", http.StatusBadRequest)
		return
	}

	// Разделяем строку с ключами проектов по запятой
	projectKeys := strings.Split(projectKeysParam, ",")

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.ReadTimeout)
	defer cancel()

	err := h.etlService.UpdateProject(ctx, projectKeys)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *Handler) GetProjects(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры запроса
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page <= 0 {
		page = 1 // Значение по умолчанию
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 20 // Значение по умолчанию
	}

	search := strings.TrimSpace(r.URL.Query().Get("search"))

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.WriteTimeout)
	defer cancel()

	// Получаем данные из JiraDB через ETL-сервис
	projects, pageInfo, err := h.etlService.GetProjectsFromJira(ctx, page, limit, search)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Формируем ответ
	response := struct {
		Projects []models.DBProject `json:"projects"`
		PageInfo models.PageInfo    `json:"pageInfo"`
	}{
		Projects: projects,
		PageInfo: pageInfo,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
