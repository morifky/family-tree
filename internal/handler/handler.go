package handler

import (
	"github.com/gin-gonic/gin"

	"brayat/internal/service"
)

// Handlers bundles all route handlers.
type Handlers struct {
	Session *SessionHandler
	// Add Person and Relationship handlers here down the line
}

// NewHandlers links all the handlers to their underlying services.
func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		Session: NewSessionHandler(services.Session),
	}
}

// RegisterRoutes registers all the routes for the application.
func (h *Handlers) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		// Sessions
		sessions := api.Group("/sessions")
		{
			sessions.POST("", h.Session.CreateSession)
			sessions.GET("/:id", h.Session.GetSession)
			sessions.PUT("/:id/status", h.Session.UpdateStatus)
			sessions.POST("/:id/extend", h.Session.ExtendExpiry)
			sessions.POST("/:id/links", h.Session.CreateAccessLink)
		}
	}
}
