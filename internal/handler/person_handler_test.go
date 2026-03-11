package handler_test

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
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

// MockPersonService is a mock of PersonService
type MockPersonService struct {
	mock.Mock
}

func (m *MockPersonService) CreatePerson(ctx context.Context, sessionID string, name string, nickname *string, gender model.Gender, photoPath *string) (*model.Person, error) {
	args := m.Called(ctx, sessionID, name, nickname, gender, photoPath)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Person), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPersonService) GetPersonByID(ctx context.Context, id string) (*model.Person, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Person), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPersonService) GetPeopleBySessionID(ctx context.Context, sessionID string) ([]model.Person, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) != nil {
		return args.Get(0).([]model.Person), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPersonService) UpdatePerson(ctx context.Context, id string, name string, nickname *string, gender model.Gender, photoPath *string) error {
	args := m.Called(ctx, id, name, nickname, gender, photoPath)
	return args.Error(0)
}

func (m *MockPersonService) DeletePerson(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockPhotoStorage is a mock of PhotoStorage
type MockPhotoStorage struct {
	mock.Mock
}

func (m *MockPhotoStorage) SavePhoto(fileHeader *multipart.FileHeader) (string, error) {
	args := m.Called(fileHeader)
	return args.String(0), args.Error(1)
}

func (m *MockPhotoStorage) DeletePhoto(filename string) error {
	args := m.Called(filename)
	return args.Error(0)
}

func setupPersonRouter(svc *MockPersonService, photo *MockPhotoStorage) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	
	h := handler.NewPersonHandler(svc, photo, zap.NewNop())
	
	router.POST("/sessions/:id/people", h.CreatePerson)
	router.GET("/sessions/:id/people", h.GetPeople)
	router.PUT("/people/:id", h.UpdatePerson)
	router.DELETE("/people/:id", h.DeletePerson)
	
	return router
}

func TestPersonHandler_CreatePerson(t *testing.T) {
	mockSvc := new(MockPersonService)
	mockPhoto := new(MockPhotoStorage)
	router := setupPersonRouter(mockSvc, mockPhoto)

	t.Run("Success No Photo", func(t *testing.T) {
		sessionID := "sess123"
		
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("name", "John")
		_ = writer.WriteField("gender", "male")
		writer.Close()

		req, _ := http.NewRequest(http.MethodPost, "/sessions/"+sessionID+"/people", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		expectedPerson := &model.Person{
			ID: "p1", Name: "John", Gender: model.GenderMale, SessionID: sessionID,
		}
		
		mockSvc.On("CreatePerson", mock.Anything, sessionID, "John", mock.Anything, model.GenderMale, mock.Anything).Return(expectedPerson, nil).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Missing Name", func(t *testing.T) {
		sessionID := "sess123"
		
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("gender", "male")
		writer.Close()

		req, _ := http.NewRequest(http.MethodPost, "/sessions/"+sessionID+"/people", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Missing Gender", func(t *testing.T) {
		sessionID := "sess123"
		
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("name", "John")
		writer.Close()

		req, _ := http.NewRequest(http.MethodPost, "/sessions/"+sessionID+"/people", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Multipart", func(t *testing.T) {
		sessionID := "sess123"
		
		req, _ := http.NewRequest(http.MethodPost, "/sessions/"+sessionID+"/people", bytes.NewReader([]byte("plain text")))
		// Set header to multipart but body is not correctly formatted multipart
		req.Header.Set("Content-Type", "multipart/form-data; boundary=myboundary")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	
	t.Run("Service Error", func(t *testing.T) {
		sessionID := "sessErr"
		
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("name", "John")
		_ = writer.WriteField("gender", "male")
		writer.Close()

		req, _ := http.NewRequest(http.MethodPost, "/sessions/"+sessionID+"/people", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		mockSvc.On("CreatePerson", mock.Anything, sessionID, "John", mock.Anything, model.GenderMale, mock.Anything).Return(nil, errors.New("svc err")).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestPersonHandler_UpdatePerson(t *testing.T) {
	mockSvc := new(MockPersonService)
	mockPhoto := new(MockPhotoStorage)
	router := setupPersonRouter(mockSvc, mockPhoto)

	t.Run("Success No Photo", func(t *testing.T) {
		pID := "p123"
		
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("name", "Jane")
		_ = writer.WriteField("gender", "female")
		writer.Close()

		req, _ := http.NewRequest(http.MethodPut, "/people/"+pID, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		mockSvc.On("UpdatePerson", mock.Anything, pID, "Jane", mock.Anything, model.GenderFemale, mock.Anything).Return(nil).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Success Remove Photo", func(t *testing.T) {
		pID := "p123"
		
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("name", "Jane")
		_ = writer.WriteField("gender", "female")
		_ = writer.WriteField("remove_photo", "true")
		writer.Close()

		req, _ := http.NewRequest(http.MethodPut, "/people/"+pID, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		mockSvc.On("UpdatePerson", mock.Anything, pID, "Jane", mock.Anything, model.GenderFemale, mock.Anything).Return(nil).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
	
	t.Run("Service Error", func(t *testing.T) {
		pID := "pErr"
		
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("name", "Jane")
		_ = writer.WriteField("gender", "female")
		writer.Close()

		req, _ := http.NewRequest(http.MethodPut, "/people/"+pID, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		mockSvc.On("UpdatePerson", mock.Anything, pID, "Jane", mock.Anything, model.GenderFemale, mock.Anything).Return(errors.New("db err")).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Missing Name", func(t *testing.T) {
		pID := "p123"
		
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("gender", "female")
		writer.Close()

		req, _ := http.NewRequest(http.MethodPut, "/people/"+pID, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Missing Gender", func(t *testing.T) {
		pID := "p123"
		
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("name", "Jane")
		writer.Close()

		req, _ := http.NewRequest(http.MethodPut, "/people/"+pID, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Multipart", func(t *testing.T) {
		pID := "p123"
		
		req, _ := http.NewRequest(http.MethodPut, "/people/"+pID, bytes.NewReader([]byte("plain text")))
		req.Header.Set("Content-Type", "multipart/form-data; boundary=myboundary")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPersonHandler_GetPeople(t *testing.T) {
	mockSvc := new(MockPersonService)
	router := setupPersonRouter(mockSvc, nil)

	t.Run("Success", func(t *testing.T) {
		sessionID := "sess123"

		req, _ := http.NewRequest(http.MethodGet, "/sessions/"+sessionID+"/people", nil)
		w := httptest.NewRecorder()

		expectedPeople := []model.Person{{ID: "p1", Name: "John"}}
		mockSvc.On("GetPeopleBySessionID", mock.Anything, sessionID).Return(expectedPeople, nil).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		sessionID := "sessErr"

		req, _ := http.NewRequest(http.MethodGet, "/sessions/"+sessionID+"/people", nil)
		w := httptest.NewRecorder()

		mockSvc.On("GetPeopleBySessionID", mock.Anything, sessionID).Return(nil, errors.New("svc err")).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestPersonHandler_DeletePerson(t *testing.T) {
	mockSvc := new(MockPersonService)
	router := setupPersonRouter(mockSvc, nil)

	t.Run("Success", func(t *testing.T) {
		pID := "p123"

		req, _ := http.NewRequest(http.MethodDelete, "/people/"+pID, nil)
		w := httptest.NewRecorder()

		mockSvc.On("DeletePerson", mock.Anything, pID).Return(nil).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		pID := "pErr"

		req, _ := http.NewRequest(http.MethodDelete, "/people/"+pID, nil)
		w := httptest.NewRecorder()

		mockSvc.On("DeletePerson", mock.Anything, pID).Return(errors.New("db err")).Once()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
