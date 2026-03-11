package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"brayat/internal/model"
)

// SessionHandler handles HTTP requests for Session resources
type SessionHandler struct {
	service model.SessionService
	logger  *zap.Logger
}

// NewSessionHandler creates a new handler serving session requests
func NewSessionHandler(svc model.SessionService, logger *zap.Logger) *SessionHandler {
	return &SessionHandler{service: svc, logger: logger}
}

// POST /sessions
func (h *SessionHandler) CreateSession(c *gin.Context) {
	var input struct {
		Title string `json:"title" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	session, err := h.service.CreateSession(c.Request.Context(), input.Title)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to create session")
		return
	}

	CreatedResponse(c, session)
}

// GET /sessions/:id
func (h *SessionHandler) GetSession(c *gin.Context) {
	id := c.Param("id")
	session, err := h.service.GetSessionByID(c.Request.Context(), id)
	if err != nil {
		ErrorResponse(c, http.StatusNotFound, "Session not found")
		return
	}

	SuccessResponse(c, session)
}

// PUT /sessions/:id/status
func (h *SessionHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := h.service.UpdateSessionStatus(c.Request.Context(), id, model.SessionStatus(input.Status))
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	SuccessResponse(c, gin.H{"message": "Status updated successfully"})
}

// POST /sessions/:id/extend
func (h *SessionHandler) ExtendExpiry(c *gin.Context) {
	id := c.Param("id")

	err := h.service.ExtendSessionExpiry(c.Request.Context(), id)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to extend session expiry")
		return
	}

	SuccessResponse(c, gin.H{"message": "Session expiry extended successfully"})
}

// POST /sessions/:id/links
func (h *SessionHandler) CreateAccessLink(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Type string `json:"type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	link, err := h.service.CreateAccessLink(c.Request.Context(), id, model.AccessType(input.Type))
	if err != nil {
		ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	CreatedResponse(c, link)
}

// GET /sessions/:id/links
func (h *SessionHandler) GetAccessLinks(c *gin.Context) {
	id := c.Param("id")
	links, err := h.service.GetAccessLinksBySessionID(c.Request.Context(), id)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to get access links")
		return
	}

	SuccessResponse(c, links)
}

// GET /sessions/verify/:code
func (h *SessionHandler) VerifyCode(c *gin.Context) {
	code := c.Param("code")
	sessionID, accessType, err := h.service.VerifyAccessCode(c.Request.Context(), code)
	if err != nil {
		ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired code")
		return
	}

	SuccessResponse(c, gin.H{
		"session_id":  sessionID,
		"access_type": accessType,
	})
}
