package app

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"jiraAnalyzer/jiraConnector/cmd/service/internal/config"
	handler "jiraAnalyzer/jiraConnector/internal/handler/http"
	"jiraAnalyzer/jiraConnector/internal/jiraclient"
	"jiraAnalyzer/jiraConnector/internal/repository/database"
	"jiraAnalyzer/jiraConnector/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type app struct {
	httpServer *http.Server
}

func NewApp(cfg config.Config) (*app, *sqlx.DB, error) {
	log.Printf("create new database config")
	db, err := database.NewDBConfig(cfg.DB)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create database config: %w", err)
	}

	clientJira := jiraclient.NewJiraClient(cfg.Jira)

	log.Printf("create new database repository")
	dbRepository := database.NewRepository(db)

	jiraService := service.NewJiraService(clientJira, dbRepository)
	etl := service.NewETLService(jiraService, dbRepository, cfg.Jira.ThreadCount, cfg.Jira.IssueInOneRequest)

	log.Printf("create new http server")
	r := mux.NewRouter()

	r.Use(handler.LogMiddleware)
	newHandler := handler.NewHandler(etl, r, cfg.Server)

	server := &http.Server{
		Addr:         cfg.Jira.Host + ":" + cfg.Server.Port,
		Handler:      newHandler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return &app{
		httpServer: server,
	}, db, nil
}

func (s *app) Run() error {
	log.Printf("Starting HTTP server on address: %s", s.httpServer.Addr)

	// Запуск HTTP-сервера в отдельной горутине
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	return nil
}

func (s *app) Close() error {
	return s.httpServer.Close()
}
