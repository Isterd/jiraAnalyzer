package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"jiraAnalyzer/backend/internal/models"
	"jiraAnalyzer/backend/internal/service"
	"jiraAnalyzer/backend/internal/utils"
	"net/http"
	"strconv"
)

type ProjectController struct {
	service *service.ProjectService
}

func NewProjectController(service *service.ProjectService) *ProjectController {
	return &ProjectController{service: service}
}

// Создание нового проекта
func (h *ProjectController) CreateProject(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("error while decoding project: %w", err))
		return
	}

	key, err := h.service.CreateProject(r.Context(), project)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("error while creating project: %w", err))
		return
	}

	project.Key = key
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

// Получение проекта по ключу
func (h *ProjectController) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("error while parse id, invalid id: %w", err))
		return
	}

	project, err := h.service.GetProjectByID(r.Context(), id)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("error while getting project by key: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(project)
}

// Получение всех проектов с пагинацией и фильтрацией
func (h *ProjectController) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")

	var page, limit int
	var err error

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid 'page' parameter: must be a positive integer"))
			return
		}
	} else {
		page = 1
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid 'limit' parameter: must be a positive integer"))
			return
		}
	} else {
		limit = 20
	}

	projects, pageInfo, err := h.service.GetAllProjects(r.Context(), page, limit, search)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("error while get all projects: %w", err))
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

// Удаление проекта по ключу
func (h *ProjectController) DeleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("error while parse id, invalid id: %w", err))
		return
	}

	if err := h.service.DeleteProject(r.Context(), id); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("error while deleting project: %w", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Обновление проекта по ключу
func (h *ProjectController) UpdateProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("error while parse id, invalid id: %w", err))
		return
	}

	var project models.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("error while decoding project: %w", err))
		return
	}

	if err := h.service.UpdateProject(r.Context(), id, project); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("error while updating project: %w", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
