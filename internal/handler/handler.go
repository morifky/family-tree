package handler

import (
	"github.com/gin-gonic/gin"

	"brayat/internal/service"
	"brayat/internal/storage"
)

// Handlers bundles all route handlers.
type Handlers struct {
	Session      *SessionHandler
	Person       *PersonHandler
	Relationship *RelationshipHandler
	Tree         *TreeHandler
}

// NewHandlers links all the handlers to their underlying services.
func NewHandlers(services *service.Services, photoStorage storage.PhotoStorage) *Handlers {
	return &Handlers{
		Session:      NewSessionHandler(services.Session),
		Person:       NewPersonHandler(services.Person, photoStorage),
		Relationship: NewRelationshipHandler(services.Relationship),
		Tree:         NewTreeHandler(services.Session, services.Person, services.Relationship),
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
}
