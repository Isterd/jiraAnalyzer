package app

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"jiraAnalyzer/backend/internal/config"
	"jiraAnalyzer/backend/internal/controller"
	"jiraAnalyzer/backend/internal/handler"
	"jiraAnalyzer/backend/internal/repository"
	"jiraAnalyzer/backend/internal/repository/database"
	"jiraAnalyzer/backend/internal/service"
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
	// Настройка логирования
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	file, err := os.OpenFile(cfg.Logging.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.Out = file
	} else {
		logger.Warnf("Failed to log to file, using default stderr: %v", err)
	}

	errFile, err := os.OpenFile(cfg.Logging.ErrorLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		errorLogger := logrus.New()
		errorLogger.SetLevel(logrus.ErrorLevel)
		errorLogger.Out = errFile
		logger.AddHook(&ErrorLogHook{Logger: errorLogger})
	} else {
		logger.Warnf("Failed to log errors to file, using default stderr: %v", err)
	}

	log.Printf("create new database config")
	db, err := database.NewDBConfig(cfg.DBSettings)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create database config: %w", err)
	}

	log.Printf("create new database repository")
	dbRepository := repository.NewRepository(db, cfg.Backend.BaseUrl)

	jiraService := service.NewService(dbRepository)

	controllers := controller.NewController(jiraService, logger, cfg.Backend)

	log.Printf("create new http server")
	r := mux.NewRouter()

	// Настройка CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200"}, // Разрешаем запросы с фронтенда
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	r.Use(c.Handler) // Применяем middleware CORS
	r.Use(handler.LogMiddleware)

	newHandler := handler.NewHandler(controllers, r)

	server := &http.Server{
		Addr:    cfg.Backend.Host + ":" + cfg.Backend.Port,
		Handler: newHandler,
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
		log.Fatalf("Backend forced to shutdown: %v", err)
	}

	return nil
}

func (s *app) Close() error {
	return s.httpServer.Close()
}

// ErrorLogHook для разделения логов ошибок
type ErrorLogHook struct {
	Logger *logrus.Logger
}

func (h *ErrorLogHook) Fire(entry *logrus.Entry) error {
	h.Logger.WithFields(entry.Data).Log(entry.Level, entry.Message)
	return nil
}

func (h *ErrorLogHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel}
}
