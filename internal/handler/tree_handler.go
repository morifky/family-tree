package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"brayat/internal/model"
)

// TreeHandler provides endpoints for returning aggregated session data (Person + Relationship)
type TreeHandler struct {
	sessionSvc      model.SessionService
	personSvc       model.PersonService
	relationshipSvc model.RelationshipService
	logger          *zap.Logger
}

// NewTreeHandler creates a new TreeHandler
func NewTreeHandler(sessionSvc model.SessionService, personSvc model.PersonService, relationshipSvc model.RelationshipService, logger *zap.Logger) *TreeHandler {
	return &TreeHandler{
		sessionSvc:      sessionSvc,
		personSvc:       personSvc,
		relationshipSvc: relationshipSvc,
		logger:          logger,
	}
}

// GetTree returns people and relationships for a session and supports ETag caching
func (h *TreeHandler) GetTree(c *gin.Context) {
	sessionID := c.Param("id")

	// Fetch Session just to check logic and Modified date
	session, err := h.sessionSvc.GetSessionByID(c.Request.Context(), sessionID)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, "Session not found")
		return
	}

	// Compute ETag from updated_at
	hash := sha256.Sum256([]byte(session.UpdatedAt.UTC().String()))
	etag := fmt.Sprintf(`"%s"`, hex.EncodeToString(hash[:]))

	// Support 304 Not Modified
	if match := c.GetHeader("If-None-Match"); match == etag {
		c.Status(http.StatusNotModified)
		return
	}

	// Has changes (or missed check), so we fetch full tree
	people, err := h.personSvc.GetPeopleBySessionID(c.Request.Context(), sessionID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch people")
		return
	}
	
	// Default to empty array instead of nil in response if no people exist
	if people == nil {
		people = []model.Person{}
	}

	relationships, err := h.relationshipSvc.GetRelationshipsBySessionID(c.Request.Context(), sessionID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch relationships")
		return
	}
	
	// Default to empty array instead of nil in response if no relationships exist
	if relationships == nil {
		relationships = []model.Relationship{}
	}

	c.Header("ETag", etag)
	SuccessResponse(c, gin.H{
		"people":        people,
		"relationships": relationships,
	})
}
