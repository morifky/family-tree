package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"brayat/internal/handler"
	"brayat/internal/model"
)

// MockSessionService is a mock of SessionService for middleware usage
type MockMiddlewareSessionService struct {
	mock.Mock
}

func (m *MockMiddlewareSessionService) CreateSession(ctx context.Context, title string) (*model.Session, error) {
	args := m.Called(ctx, title)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Session), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMiddlewareSessionService) GetSessionByID(ctx context.Context, id string) (*model.Session, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Session), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMiddlewareSessionService) GetSessionByAdminCode(ctx context.Context, code string) (*model.Session, error) {
	args := m.Called(ctx, code)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Session), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMiddlewareSessionService) UpdateSessionStatus(ctx context.Context, id string, status model.SessionStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockMiddlewareSessionService) ExtendSessionExpiry(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMiddlewareSessionService) CreateAccessLink(ctx context.Context, sessionID string, accessType model.AccessType) (*model.AccessLink, error) {
	args := m.Called(ctx, sessionID, accessType)
	if args.Get(0) != nil {
		return args.Get(0).(*model.AccessLink), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMiddlewareSessionService) VerifyAccessCode(ctx context.Context, code string) (string, model.AccessType, error) {
	args := m.Called(ctx, code)
	return args.String(0), args.Get(1).(model.AccessType), args.Error(2)
}

func (m *MockMiddlewareSessionService) GetAccessLinksBySessionID(ctx context.Context, sessionID string) ([]model.AccessLink, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) != nil {
		return args.Get(0).([]model.AccessLink), args.Error(1)
	}
	return nil, args.Error(1)
}

func setupMiddlewareRouter(svc model.SessionService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	auth := router.Group("/api", handler.AuthMiddleware(svc))
	{
		auth.GET("/protected", func(c *gin.Context) {
			sessionID, _ := c.Get(handler.ContextSessionIDKey)
			access, _ := c.Get(handler.ContextAccessLevelKey)

			c.JSON(http.StatusOK, gin.H{
				"session_id": sessionID,
				"access":     access,
			})
		})
	}
	return router
}

func TestAuthMiddleware(t *testing.T) {
	mockSvc := new(MockMiddlewareSessionService)
	router := setupMiddlewareRouter(mockSvc)

	t.Run("Missing Authorization Header", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/protected", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Missing Authorization header")
	})

	t.Run("Invalid Authorization Format - Not Bearer", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/protected", nil)
		req.Header.Set("Authorization", "Basic token123")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid Authorization header format")
	})

	t.Run("Invalid Authorization Format - Missing Token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/protected", nil)
		req.Header.Set("Authorization", "Bearer")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid Authorization header format")
	})

	t.Run("Invalid Access Code", func(t *testing.T) {
		mockSvc.On("VerifyAccessCode", mock.Anything, "invalid_code").Return("", model.AccessType(""), errors.New("invalid code")).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/protected", nil)
		req.Header.Set("Authorization", "Bearer invalid_code")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid access code")
		mockSvc.AssertExpectations(t)
	})

	t.Run("Valid Session, Active Status", func(t *testing.T) {
		mockSvc.On("VerifyAccessCode", mock.Anything, "valid_code").Return("sess123", model.AccessTypeAdmin, nil).Once()
		mockSvc.On("GetSessionByID", mock.Anything, "sess123").Return(&model.Session{ID: "sess123", Status: model.SessionStatusActive}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/protected", nil)
		req.Header.Set("Authorization", "Bearer valid_code")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"session_id":"sess123"`)
		assert.Contains(t, w.Body.String(), `"access":"admin"`)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Valid Code, Session Not Found", func(t *testing.T) {
		mockSvc.On("VerifyAccessCode", mock.Anything, "valid_code2").Return("sess_miss", model.AccessTypeView, nil).Once()
		mockSvc.On("GetSessionByID", mock.Anything, "sess_miss").Return(nil, errors.New("not found")).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/protected", nil)
		req.Header.Set("Authorization", "Bearer valid_code2")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Session not found")
		mockSvc.AssertExpectations(t)
	})

	t.Run("Valid Code, Locked Session", func(t *testing.T) {
		mockSvc.On("VerifyAccessCode", mock.Anything, "valid_code3").Return("sess_locked", model.AccessTypeView, nil).Once()
		mockSvc.On("GetSessionByID", mock.Anything, "sess_locked").Return(&model.Session{ID: "sess_locked", Status: model.SessionStatusLocked}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/protected", nil)
		req.Header.Set("Authorization", "Bearer valid_code3")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "Session is locked")
		mockSvc.AssertExpectations(t)
	})

	t.Run("Valid Code, Closed Session", func(t *testing.T) {
		mockSvc.On("VerifyAccessCode", mock.Anything, "valid_code4").Return("sess_closed", model.AccessTypeView, nil).Once()
		mockSvc.On("GetSessionByID", mock.Anything, "sess_closed").Return(&model.Session{ID: "sess_closed", Status: model.SessionStatusClosed}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/protected", nil)
		req.Header.Set("Authorization", "Bearer valid_code4")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "Session is closed")
		mockSvc.AssertExpectations(t)
	})
}
