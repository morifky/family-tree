package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"brayat/internal/model"
)

// RelationshipHandler handles HTTP requests for Relationship resources
type RelationshipHandler struct {
	service model.RelationshipService
	logger  *zap.Logger
}

// NewRelationshipHandler creates a new handler serving relationship requests
func NewRelationshipHandler(svc model.RelationshipService, logger *zap.Logger) *RelationshipHandler {
	return &RelationshipHandler{service: svc, logger: logger}
}

// POST /sessions/:id/relationships
func (h *RelationshipHandler) CreateRelationship(c *gin.Context) {
	sessionID := c.Param("id")

	var input struct {
		PersonAID string                 `json:"person_a_id" binding:"required"`
		PersonBID string                 `json:"person_b_id" binding:"required"`
		Type      model.RelationshipType `json:"type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	rel, err := h.service.CreateRelationship(c.Request.Context(), sessionID, input.PersonAID, input.PersonBID, input.Type)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create relationship")
		return
	}

	CreatedResponse(c, rel)
}

// GET /sessions/:id/relationships
func (h *RelationshipHandler) GetRelationships(c *gin.Context) {
	sessionID := c.Param("id")

	rels, err := h.service.GetRelationshipsBySessionID(c.Request.Context(), sessionID)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve relationships")
		return
	}

	SuccessResponse(c, rels)
}

// DELETE /relationships/:id
func (h *RelationshipHandler) DeleteRelationship(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeleteRelationship(c.Request.Context(), id)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to delete relationship")
		return
	}

	SuccessResponse(c, gin.H{"message": "Relationship deleted successfully"})
}
