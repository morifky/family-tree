package handler_test

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"brayat/internal/handler"
	"brayat/internal/model"
)

func setupTreeRouter(
	sessionSvc *MockSessionService,
	personSvc *MockPersonService,
	relSvc *MockRelationshipService,
) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	treeHandler := handler.NewTreeHandler(sessionSvc, personSvc, relSvc, zap.NewNop())
	router.GET("/sessions/:id/tree", treeHandler.GetTree)

	return router
}

func TestTreeHandler_GetTree(t *testing.T) {
	sessionSvc := new(MockSessionService)
	personSvc := new(MockPersonService)
	relSvc := new(MockRelationshipService)

	router := setupTreeRouter(sessionSvc, personSvc, relSvc)

	t.Run("Success full fetch", func(t *testing.T) {
		sessionID := "sess123"
		updatedAt := time.Now()

		session := &model.Session{
			ID:        sessionID,
			UpdatedAt: updatedAt,
		}

		people := []model.Person{
			{ID: "p1", Name: "John Doe"},
			{ID: "p2", Name: "Jane Doe"},
		}

		rels := []model.Relationship{
			{ID: "r1", PersonAID: "p1", PersonBID: "p2", Type: model.RelationshipTypeSpouse},
		}

		sessionSvc.On("GetSessionByID", mock.Anything, sessionID).Return(session, nil).Once()
		personSvc.On("GetPeopleBySessionID", mock.Anything, sessionID).Return(people, nil).Once()
		relSvc.On("GetRelationshipsBySessionID", mock.Anything, sessionID).Return(rels, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/sessions/"+sessionID+"/tree", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Validate ETag Headers
		hash := sha256.Sum256([]byte(updatedAt.UTC().String()))
		expectedETag := fmt.Sprintf(`"%s"`, hex.EncodeToString(hash[:]))
		assert.Equal(t, expectedETag, w.Header().Get("ETag"))

		sessionSvc.AssertExpectations(t)
		personSvc.AssertExpectations(t)
		relSvc.AssertExpectations(t)
	})

	t.Run("304 Not Modified", func(t *testing.T) {
		sessionID := "sess123"
		updatedAt := time.Now()

		session := &model.Session{
			ID:        sessionID,
			UpdatedAt: updatedAt,
		}

		sessionSvc.On("GetSessionByID", mock.Anything, sessionID).Return(session, nil).Once()

		// Calculate matching ETag
		hash := sha256.Sum256([]byte(updatedAt.UTC().String()))
		etag := fmt.Sprintf(`"%s"`, hex.EncodeToString(hash[:]))

		req, _ := http.NewRequest(http.MethodGet, "/sessions/"+sessionID+"/tree", nil)
		req.Header.Set("If-None-Match", etag) // Provide matching ETag
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotModified, w.Code)
		assert.Empty(t, w.Body.String()) // No body should be sent

		// Assert zero calls to fetch data since ETag matched
		sessionSvc.AssertExpectations(t)
	})

	t.Run("Session Not Found", func(t *testing.T) {
		sessionID := "missing"

		sessionSvc.On("GetSessionByID", mock.Anything, sessionID).Return(nil, errors.New("not found")).Once()

		req, _ := http.NewRequest(http.MethodGet, "/sessions/"+sessionID+"/tree", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		sessionSvc.AssertExpectations(t)
	})

	t.Run("People Fetch Error", func(t *testing.T) {
		sessionID := "sess123"
		updatedAt := time.Now()

		session := &model.Session{
			ID:        sessionID,
			UpdatedAt: updatedAt,
		}

		sessionSvc.On("GetSessionByID", mock.Anything, sessionID).Return(session, nil).Once()
		personSvc.On("GetPeopleBySessionID", mock.Anything, sessionID).Return(nil, errors.New("db error")).Once()

		req, _ := http.NewRequest(http.MethodGet, "/sessions/"+sessionID+"/tree", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		sessionSvc.AssertExpectations(t)
		personSvc.AssertExpectations(t)
	})

	t.Run("Relationships Fetch Error", func(t *testing.T) {
		sessionID := "sess123"
		updatedAt := time.Now()

		session := &model.Session{
			ID:        sessionID,
			UpdatedAt: updatedAt,
		}

		people := []model.Person{
			{ID: "p1", Name: "John Doe"},
		}

		sessionSvc.On("GetSessionByID", mock.Anything, sessionID).Return(session, nil).Once()
		personSvc.On("GetPeopleBySessionID", mock.Anything, sessionID).Return(people, nil).Once()
		relSvc.On("GetRelationshipsBySessionID", mock.Anything, sessionID).Return(nil, errors.New("db error")).Once()

		req, _ := http.NewRequest(http.MethodGet, "/sessions/"+sessionID+"/tree", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		sessionSvc.AssertExpectations(t)
		personSvc.AssertExpectations(t)
		relSvc.AssertExpectations(t)
	})
}
