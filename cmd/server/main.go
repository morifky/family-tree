package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"brayat/internal/config"
	"brayat/internal/db"
	"brayat/internal/handler"
	"brayat/internal/repository"
	"brayat/internal/service"
	"brayat/internal/storage"
)

func main() {
	// Initialize Logger
	logger, _ := zap.NewProduction() // Use production logger by default
	defer logger.Sync()

	logger.Info("Starting Brayat Server...")

	// Load Config
	cfg := config.MustLoad()

	// Connect to Database
	gormDB := db.MustOpen(cfg.DatabasePath)

	// Keep an instance of sqlDB to close on shutdown
	sqlDB, err := gormDB.DB()
	if err != nil {
		logger.Fatal("Failed to get sql.DB from gorm", zap.Error(err))
	}

	// Initialize storage
	photoStorage := storage.NewLocalPhotoStorage(cfg.PhotosDir)

	// Initialize repositories
	repos := &repository.Repositories{
		Session:      repository.NewSessionRepository(gormDB),
		Person:       repository.NewPersonRepository(gormDB),
		Relationship: repository.NewRelationshipRepository(gormDB),
	}

	// Initialize services
	services := &service.Services{
		Session:      service.NewSessionService(repos.Session),
		Person:       service.NewPersonService(repos.Person, photoStorage),
		Relationship: service.NewRelationshipService(repos.Relationship),
	}

	// Initialize handlers with logger
	handlers := handler.NewHandlers(services, photoStorage, logger)

	// Setup Router
	// Consider adding Gin Zap Middleware here if needed, but Gin's default is fine for now
	router := gin.Default()
	handlers.RegisterRoutes(router)

	// Configure HTTP Server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	// Start server in background
	go func() {
		logger.Info("Listening", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Failed to bind / start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// The context informs the server it has 5 seconds to finish the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", zap.Error(err))
	}

	logger.Info("Closing database connection...")
	if err := sqlDB.Close(); err != nil {
		logger.Error("Failed to close database cleanly", zap.Error(err))
	}

	logger.Info("Server exiting")
}
