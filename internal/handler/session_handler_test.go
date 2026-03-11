package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"brayat/internal/handler"
	"brayat/internal/model"
)

// MockSessionService is a mock of SessionService
type MockSessionService struct {
	mock.Mock
}

func (m *MockSessionService) CreateSession(ctx context.Context, title string) (*model.Session, error) {
	args := m.Called(ctx, title)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Session), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSessionService) GetSessionByID(ctx context.Context, id string) (*model.Session, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Session), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSessionService) GetSessionByAdminCode(ctx context.Context, code string) (*model.Session, error) {
	args := m.Called(ctx, code)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Session), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSessionService) UpdateSessionStatus(ctx context.Context, id string, status model.SessionStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockSessionService) ExtendSessionExpiry(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSessionService) CreateAccessLink(ctx context.Context, sessionID string, accessType model.AccessType) (*model.AccessLink, error) {
	args := m.Called(ctx, sessionID, accessType)
	if args.Get(0) != nil {
		return args.Get(0).(*model.AccessLink), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSessionService) VerifyAccessCode(ctx context.Context, code string) (string, model.AccessType, error) {
	args := m.Called(ctx, code)
	return args.String(0), args.Get(1).(model.AccessType), args.Error(2)
}

func (m *MockSessionService) GetAccessLinksBySessionID(ctx context.Context, sessionID string) ([]model.AccessLink, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) != nil {
		return args.Get(0).([]model.AccessLink), args.Error(1)
	}
	return nil, args.Error(1)
}

func setupSessionRouter(svc *MockSessionService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	
	h := handler.NewSessionHandler(svc, zap.NewNop())
	
	router.POST("/sessions", h.CreateSession)
	router.GET("/sessions/:id", h.GetSession)
	router.PUT("/sessions/:id/status", h.UpdateStatus)
	router.POST("/sessions/:id/extend", h.ExtendExpiry)
	router.POST("/sessions/:id/links", h.CreateAccessLink)
	
	return router
}

func TestSessionHandler_CreateSession(t *testing.T) {
	mockSvc := new(MockSessionService)
	router := setupSessionRouter(mockSvc)

	t.Run("Success", func(t *testing.T) {
		input := map[string]string{"title": "My Family"}
		body, _ := json.Marshal(input)
		req, _ := http.NewRequest(http.MethodPost, "/sessions", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		expectedSess := &model.Session{ID: "sess1", Title: "My Family"}
		mockSvc.On("CreateSession", mock.Anything, "My Family").Return(expectedSess, nil).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/sessions", bytes.NewBuffer([]byte("{invalid}")))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	
	t.Run("Service Error", func(t *testing.T) {
		input := map[string]string{"title": "My Family"}
		body, _ := json.Marshal(input)
		req, _ := http.NewRequest(http.MethodPost, "/sessions", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		mockSvc.On("CreateSession", mock.Anything, "My Family").Return(nil, errors.New("err")).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestSessionHandler_GetSession(t *testing.T) {
	mockSvc := new(MockSessionService)
	router := setupSessionRouter(mockSvc)

	t.Run("Success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/sessions/sess1", nil)
		w := httptest.NewRecorder()

		s := &model.Session{ID: "sess1", Title: "t"}
		mockSvc.On("GetSessionByID", mock.Anything, "sess1").Return(s, nil).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/sessions/sess2", nil)
		w := httptest.NewRecorder()

		mockSvc.On("GetSessionByID", mock.Anything, "sess2").Return(nil, errors.New("err")).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestSessionHandler_UpdateStatus(t *testing.T) {
	mockSvc := new(MockSessionService)
	router := setupSessionRouter(mockSvc)

	t.Run("Success", func(t *testing.T) {
		input := map[string]string{"status": string(model.SessionStatusLocked)}
		body, _ := json.Marshal(input)
		req, _ := http.NewRequest(http.MethodPut, "/sessions/sess1/status", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		mockSvc.On("UpdateSessionStatus", mock.Anything, "sess1", model.SessionStatusLocked).Return(nil).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Invalid Input", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/sessions/sess1/status", bytes.NewBuffer([]byte("{invalid}")))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		input := map[string]string{"status": "invalid_status"}
		body, _ := json.Marshal(input)
		req, _ := http.NewRequest(http.MethodPut, "/sessions/sess1/status", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		mockSvc.On("UpdateSessionStatus", mock.Anything, "sess1", model.SessionStatus("invalid_status")).Return(errors.New("invalid session status")).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestSessionHandler_ExtendExpiry(t *testing.T) {
	mockSvc := new(MockSessionService)
	router := setupSessionRouter(mockSvc)

	t.Run("Success", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/sessions/sess1/extend", nil)
		w := httptest.NewRecorder()

		mockSvc.On("ExtendSessionExpiry", mock.Anything, "sess1").Return(nil).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/sessions/sess2/extend", nil)
		w := httptest.NewRecorder()

		mockSvc.On("ExtendSessionExpiry", mock.Anything, "sess2").Return(errors.New("db err")).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestSessionHandler_CreateAccessLink(t *testing.T) {
	mockSvc := new(MockSessionService)
	router := setupSessionRouter(mockSvc)

	t.Run("Success", func(t *testing.T) {
		input := map[string]string{"type": string(model.AccessTypeEdit)}
		body, _ := json.Marshal(input)
		req, _ := http.NewRequest(http.MethodPost, "/sessions/sess1/links", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		link := &model.AccessLink{ID: "link1", Type: model.AccessTypeEdit}
		mockSvc.On("CreateAccessLink", mock.Anything, "sess1", model.AccessTypeEdit).Return(link, nil).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Invalid Input", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/sessions/sess1/links", bytes.NewBuffer([]byte("{invalid}")))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		input := map[string]string{"type": "any"}
		body, _ := json.Marshal(input)
		req, _ := http.NewRequest(http.MethodPost, "/sessions/sess2/links", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		mockSvc.On("CreateAccessLink", mock.Anything, "sess2", model.AccessType("any")).Return(nil, errors.New("err")).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
