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

type IssueController struct {
	service *service.IssueService
}

func NewIssueController(service *service.IssueService) *IssueController {
	return &IssueController{service: service}
}

func (h *IssueController) GetAllIssues(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

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

	issues, err := h.service.GetAllIssues(r.Context(), page, limit)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to get issues: %w", err))
		return
	}

	utils.WriteJSONResponse(w, map[string]interface{}{
		"issues": issues,
	})
}

func (h *IssueController) GetIssueById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("error while parse id, invalid id: %w", err))
		return
	}

	issue, err := h.service.GetIssueById(r.Context(), id)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to get issues: %w", err))
		return
	}

	utils.WriteJSONResponse(w, map[string]interface{}{
		"issue": issue,
	})
}

func (h *IssueController) CreateIssue(w http.ResponseWriter, r *http.Request) {
	var issue models.Issue
	if err := json.NewDecoder(r.Body).Decode(&issue); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	if err := h.service.CreateIssue(r.Context(), issue); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to create issue: %w", err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *IssueController) UpdateIssue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("error while parse id, invalid id: %w", err))
		return
	}

	var issue models.Issue
	if err := json.NewDecoder(r.Body).Decode(&issue); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
		return
	}

	issue.ID = id
	if err := h.service.UpdateIssue(r.Context(), issue); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to update issue: ", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *IssueController) DeleteIssue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("error while parse id, invalid id: %w", err))
		return
	}

	if err := h.service.DeleteIssue(r.Context(), id); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to delete issue: ", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
