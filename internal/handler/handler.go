package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"brayat/internal/service"
	"brayat/internal/storage"
)

// Handlers bundles all route handlers.
type Handlers struct {
	Session      *SessionHandler
	Person       *PersonHandler
	Relationship *RelationshipHandler
	Tree         *TreeHandler
	logger       *zap.Logger
}

// NewHandlers links all the handlers to their underlying services.
func NewHandlers(services *service.Services, photoStorage storage.PhotoStorage, logger *zap.Logger) *Handlers {
	return &Handlers{
		Session:      NewSessionHandler(services.Session, logger),
		Person:       NewPersonHandler(services.Person, photoStorage, logger),
		Relationship: NewRelationshipHandler(services.Relationship, logger),
		Tree:         NewTreeHandler(services.Session, services.Person, services.Relationship, logger),
		logger:       logger,
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
			sessions.GET("/verify/:code", h.Session.VerifyCode)
			sessions.GET("/:id", h.Session.GetSession)
			sessions.PUT("/:id/status", h.Session.UpdateStatus)
			sessions.POST("/:id/extend", h.Session.ExtendExpiry)
			sessions.GET("/:id/links", h.Session.GetAccessLinks)
			sessions.POST("/:id/links", h.Session.CreateAccessLink)

			// Nested People under Session
			sessions.POST("/:id/people", h.Person.CreatePerson)
			sessions.GET("/:id/people", h.Person.GetPeople)

			// Nested Relationships under Session
			sessions.POST("/:id/relationships", h.Relationship.CreateRelationship)
			sessions.GET("/:id/relationships", h.Relationship.GetRelationships)

			// Aggregate Polling Endpoints
			sessions.GET("/:id/tree", h.Tree.GetTree)
		}

		// People
		people := api.Group("/people")
		{
			people.PUT("/:id", h.Person.UpdatePerson)
			people.DELETE("/:id", h.Person.DeletePerson)
		}

		// Relationships
		relationships := api.Group("/relationships")
		{
			relationships.DELETE("/:id", h.Relationship.DeleteRelationship)
		}
	}

	// Serve Photos from /data/photos
	// In production, cfg.PhotosDir defaults to /data/photos
	router.Static("/photos", "/data/photos")

	// Serve Frontend (SvelteKit static build)
	// We serve the build directory and fallback to index.html for SPA routing
	router.Static("/_app", "./web/build/_app")
	router.StaticFile("/favicon.png", "./web/build/favicon.png")
	router.StaticFile("/robots.txt", "./web/build/robots.txt")

	router.NoRoute(func(c *gin.Context) {
		// If requesting something that looks like an asset but wasn't found, 404
		if strings.Contains(c.Request.URL.Path, ".") {
			c.Status(http.StatusNotFound)
			return
		}
		// Otherwise serve index.html for SPA routing
		c.File("./web/build/index.html")
	})
}
