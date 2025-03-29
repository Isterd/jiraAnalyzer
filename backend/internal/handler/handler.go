package handler

import (
	"github.com/gorilla/mux"
	"jiraAnalyzer/backend/internal/controller"
	"net/http"
)

func NewHandler(controllers *controller.Controller, r *mux.Router) *mux.Router {
	setProjectRoute(controllers.ProjectController, r)
	setAnalyticRoute(controllers.AnalyticsController, r)
	setIssueRoute(controllers.IssueController, r)
	setConnectorRoute(controllers.JiraController, r)

	return r
}

func setAnalyticRoute(ac *controller.AnalyticsController, r *mux.Router) {
	r.HandleFunc("/api/v1/projects/{key}/analytics", ac.GetProjectAnalytics).Methods(http.MethodOptions, http.MethodGet)
	r.HandleFunc("/api/v1/isAnalyzed/{key}", ac.IsProjectAnalytics).Methods(http.MethodOptions, http.MethodGet)
	r.HandleFunc("/api/v1/delete/{key}", ac.DeleteProjectAnalytics).Methods(http.MethodOptions, http.MethodDelete)
	r.HandleFunc("/api/v1/graph/get/{taskNumber:[0-9]+}", ac.GetGraph).Methods(http.MethodOptions, http.MethodGet)
	r.HandleFunc("/api/v1/graph/make/{taskNumber:[0-9]}", ac.MakeGraph).Methods(http.MethodOptions, http.MethodPost)
	r.HandleFunc("/api/v1/compare/{taskNumber:[0-9]+}", ac.GetComparison).Methods(http.MethodOptions, http.MethodGet)
	r.HandleFunc("/api/v1/compare/{taskNumber:[0-9]}", ac.MakeComparison).Methods(http.MethodOptions, http.MethodPost)
}

func setProjectRoute(pc *controller.ProjectController, r *mux.Router) {
	r.HandleFunc("/api/v1/projects", pc.CreateProject).Methods(http.MethodOptions, http.MethodPost)
	r.HandleFunc("/api/v1/projects/{id}", pc.DeleteProject).Methods(http.MethodOptions, http.MethodDelete)
	r.HandleFunc("/api/v1/projects/{id}", pc.GetProjectByID).Methods(http.MethodOptions, http.MethodGet)
	r.HandleFunc("/api/v1/projects", pc.GetAllProjects).Methods(http.MethodOptions, http.MethodGet)
	r.HandleFunc("/api/v1/putProject/{id}", pc.UpdateProject).Methods(http.MethodOptions, http.MethodPut)
}

func setIssueRoute(ic *controller.IssueController, r *mux.Router) {
	r.HandleFunc("/api/v1/issues", ic.GetAllIssues).Methods(http.MethodOptions, http.MethodGet)
	r.HandleFunc("/api/v1/issues/{id}", ic.GetIssueById).Methods(http.MethodOptions, http.MethodGet)
	r.HandleFunc("/api/v1/issues", ic.CreateIssue).Methods(http.MethodOptions, http.MethodPost)
	r.HandleFunc("/api/v1/issues/{id}", ic.UpdateIssue).Methods(http.MethodOptions, http.MethodPut)
	r.HandleFunc("/api/v1/issues/{id}", ic.DeleteIssue).Methods(http.MethodOptions, http.MethodDelete)
}

func setConnectorRoute(jc *controller.JiraController, r *mux.Router) {
	r.HandleFunc("/api/v1/connector/updateProject", jc.UpdateConnectorProject).Methods(http.MethodOptions, http.MethodPost)
	r.HandleFunc("/api/v1/connector/projects", jc.GetConnectorProjects).Methods(http.MethodOptions, http.MethodGet)
}
