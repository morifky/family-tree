package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"brayat/internal/model"
)

const (
	ContextSessionIDKey   = "session_id"
	ContextAccessLevelKey = "access"
)

// AuthMiddleware creates a gin middleware to enforce session and access link authentication.
func AuthMiddleware(sessionSvc model.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			ErrorResponse(c, http.StatusUnauthorized, "Missing Authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			ErrorResponse(c, http.StatusUnauthorized, "Invalid Authorization header format")
			c.Abort()
			return
		}

		code := parts[1]

		sessionID, accessType, err := sessionSvc.VerifyAccessCode(c.Request.Context(), code)
		if err != nil {
			ErrorResponse(c, http.StatusUnauthorized, "Invalid access code")
			c.Abort()
			return
		}

		// Verify session status
		session, err := sessionSvc.GetSessionByID(c.Request.Context(), sessionID)
		if err != nil {
			ErrorResponse(c, http.StatusUnauthorized, "Session not found")
			c.Abort()
			return
		}

		if session.Status == model.SessionStatusLocked || session.Status == model.SessionStatusClosed {
			ErrorResponse(c, http.StatusForbidden, "Session is "+string(session.Status))
			c.Abort()
			return
		}

		c.Set(ContextSessionIDKey, sessionID)
		c.Set(ContextAccessLevelKey, accessType)

		c.Next()
	}
}
