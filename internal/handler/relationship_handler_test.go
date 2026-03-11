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

// MockRelationshipService is a mock of RelationshipService
type MockRelationshipService struct {
	mock.Mock
}

func (m *MockRelationshipService) CreateRelationship(ctx context.Context, sessionID string, personAID string, personBID string, relType model.RelationshipType) (*model.Relationship, error) {
	args := m.Called(ctx, sessionID, personAID, personBID, relType)
	var rel *model.Relationship
	if args.Get(0) != nil {
		rel = args.Get(0).(*model.Relationship)
	}
	return rel, args.Error(1)
}

func (m *MockRelationshipService) GetRelationshipsBySessionID(ctx context.Context, sessionID string) ([]model.Relationship, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) != nil {
		return args.Get(0).([]model.Relationship), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRelationshipService) DeleteRelationship(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupRelationshipRouter(svc *MockRelationshipService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	h := handler.NewRelationshipHandler(svc, zap.NewNop())

	router.POST("/sessions/:id/relationships", h.CreateRelationship)
	router.GET("/sessions/:id/relationships", h.GetRelationships)
	router.DELETE("/relationships/:id", h.DeleteRelationship)

	return router
}

func TestRelationshipHandler_CreateRelationship(t *testing.T) {
	mockSvc := new(MockRelationshipService)
	router := setupRelationshipRouter(mockSvc)

	t.Run("Success", func(t *testing.T) {
		sessionID := "sess123"
		input := map[string]interface{}{
			"person_a_id": "pA1",
			"person_b_id": "pB2",
			"type":        model.RelationshipTypeParentChild,
		}
		body, _ := json.Marshal(input)

		req, _ := http.NewRequest(http.MethodPost, "/sessions/"+sessionID+"/relationships", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Expected response mock
		expectedRel := &model.Relationship{
			ID:        "rel1",
			SessionID: sessionID,
			PersonAID: "pA1",
			PersonBID: "pB2",
			Type:      model.RelationshipTypeParentChild,
		}
		mockSvc.On("CreateRelationship", mock.Anything, sessionID, "pA1", "pB2", model.RelationshipTypeParentChild).Return(expectedRel, nil).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Invalid Input", func(t *testing.T) {
		sessionID := "sess123"
		input := map[string]interface{}{
			// Missing person_a_id
			"person_b_id": "pB2",
			"type":        model.RelationshipTypeParentChild,
		}
		body, _ := json.Marshal(input)

		req, _ := http.NewRequest(http.MethodPost, "/sessions/"+sessionID+"/relationships", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		sessionID := "sess123"
		input := map[string]interface{}{
			"person_a_id": "pA1",
			"person_b_id": "pB2",
			"type":        model.RelationshipTypeParentChild,
		}
		body, _ := json.Marshal(input)

		req, _ := http.NewRequest(http.MethodPost, "/sessions/"+sessionID+"/relationships", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockSvc.On("CreateRelationship", mock.Anything, sessionID, "pA1", "pB2", model.RelationshipTypeParentChild).Return(nil, errors.New("svc err")).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestRelationshipHandler_GetRelationships(t *testing.T) {
	mockSvc := new(MockRelationshipService)
	router := setupRelationshipRouter(mockSvc)

	t.Run("Success", func(t *testing.T) {
		sessionID := "sess123"

		req, _ := http.NewRequest(http.MethodGet, "/sessions/"+sessionID+"/relationships", nil)
		w := httptest.NewRecorder()

		expectedRels := []model.Relationship{
			{ID: "rel1", SessionID: sessionID},
		}
		mockSvc.On("GetRelationshipsBySessionID", mock.Anything, sessionID).Return(expectedRels, nil).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		sessionID := "sess123"

		req, _ := http.NewRequest(http.MethodGet, "/sessions/"+sessionID+"/relationships", nil)
		w := httptest.NewRecorder()

		mockSvc.On("GetRelationshipsBySessionID", mock.Anything, sessionID).Return(nil, errors.New("svc err")).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestRelationshipHandler_DeleteRelationship(t *testing.T) {
	mockSvc := new(MockRelationshipService)
	router := setupRelationshipRouter(mockSvc)

	t.Run("Success", func(t *testing.T) {
		relID := "rel123"

		req, _ := http.NewRequest(http.MethodDelete, "/relationships/"+relID, nil)
		w := httptest.NewRecorder()

		mockSvc.On("DeleteRelationship", mock.Anything, relID).Return(nil).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		relID := "rel123"

		req, _ := http.NewRequest(http.MethodDelete, "/relationships/"+relID, nil)
		w := httptest.NewRecorder()

		mockSvc.On("DeleteRelationship", mock.Anything, relID).Return(errors.New("db err")).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
