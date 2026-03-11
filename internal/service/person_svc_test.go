package service_test

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"brayat/internal/model"
	"brayat/internal/service"
)

// MockPersonRepository is a mock of PersonRepository
type MockPersonRepository struct {
	mock.Mock
}

func (m *MockPersonRepository) CreatePerson(ctx context.Context, person *model.Person) error {
	args := m.Called(ctx, person)
	return args.Error(0)
}

func (m *MockPersonRepository) GetPersonByID(ctx context.Context, id string) (*model.Person, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Person), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPersonRepository) GetPeopleBySessionID(ctx context.Context, sessionID string) ([]model.Person, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) != nil {
		return args.Get(0).([]model.Person), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPersonRepository) UpdatePerson(ctx context.Context, person *model.Person) error {
	args := m.Called(ctx, person)
	return args.Error(0)
}

func (m *MockPersonRepository) DeletePerson(ctx context.Context, id string) error {
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

func TestPersonService_CreatePerson(t *testing.T) {
	mockRepo := new(MockPersonRepository)
	mockPhoto := new(MockPhotoStorage)
	svc := service.NewPersonService(mockRepo, mockPhoto)

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		sessionID := "sess1"
		name := "John"
		nickname := "Johnny"
		gender := model.GenderMale
		photoPath := "test.jpg"

		mockRepo.On("CreatePerson", ctx, mock.AnythingOfType("*model.Person")).Return(nil).Once()

		p, err := svc.CreatePerson(ctx, sessionID, name, &nickname, gender, &photoPath)
		assert.NoError(t, err)
		assert.Equal(t, name, p.Name)
		assert.Equal(t, photoPath, *p.PhotoPath)

		mockRepo.AssertExpectations(t)
	})
	
	t.Run("Error", func(t *testing.T) {
		sessionID := "sess1"
		name := "John"
		gender := model.GenderMale

		mockRepo.On("CreatePerson", ctx, mock.AnythingOfType("*model.Person")).Return(errors.New("db err")).Once()

		p, err := svc.CreatePerson(ctx, sessionID, name, nil, gender, nil)
		assert.Error(t, err)
		assert.Nil(t, p)

		mockRepo.AssertExpectations(t)
	})
}

func TestPersonService_GetPersonByID(t *testing.T) {
	mockRepo := new(MockPersonRepository)
	svc := service.NewPersonService(mockRepo, nil)
	ctx := context.Background()

	p := &model.Person{ID: "p1", Name: "Jane"}
	mockRepo.On("GetPersonByID", ctx, "p1").Return(p, nil).Once()

	res, err := svc.GetPersonByID(ctx, "p1")
	assert.NoError(t, err)
	assert.Equal(t, p, res)
	mockRepo.AssertExpectations(t)
}

func TestPersonService_GetPeopleBySessionID(t *testing.T) {
	mockRepo := new(MockPersonRepository)
	svc := service.NewPersonService(mockRepo, nil)
	ctx := context.Background()

	people := []model.Person{{ID: "p1"}, {ID: "p2"}}
	mockRepo.On("GetPeopleBySessionID", ctx, "sess1").Return(people, nil).Once()

	res, err := svc.GetPeopleBySessionID(ctx, "sess1")
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	mockRepo.AssertExpectations(t)
}

func TestPersonService_UpdatePerson(t *testing.T) {
	mockRepo := new(MockPersonRepository)
	mockPhoto := new(MockPhotoStorage)
	svc := service.NewPersonService(mockRepo, mockPhoto)
	ctx := context.Background()

	t.Run("Success - No Photo Change", func(t *testing.T) {
		oldPhoto := "old.jpg"
		p := &model.Person{ID: "p1", Name: "Old", PhotoPath: &oldPhoto}

		mockRepo.On("GetPersonByID", ctx, "p1").Return(p, nil).Once()
		mockRepo.On("UpdatePerson", ctx, mock.AnythingOfType("*model.Person")).Return(nil).Once()

		// Not modifying photo
		err := svc.UpdatePerson(ctx, "p1", "New", nil, model.GenderFemale, nil)
		assert.NoError(t, err)
		assert.Equal(t, "New", p.Name)
		assert.Equal(t, "old.jpg", *p.PhotoPath) // Keeps old since nil passed

		mockRepo.AssertExpectations(t)
		mockPhoto.AssertExpectations(t)
	})

	t.Run("Success - Remove Photo", func(t *testing.T) {
		// Mock setup
		oldPhoto := "old.jpg"
		p := &model.Person{ID: "p2", Name: "Old", PhotoPath: &oldPhoto}
		
		mockRepo.On("GetPersonByID", ctx, "p2").Return(p, nil).Once()
		mockPhoto.On("DeletePhoto", "old.jpg").Return(nil).Once()
		mockRepo.On("UpdatePerson", ctx, mock.AnythingOfType("*model.Person")).Return(nil).Once()

		empty := "" // signals delete
		err := svc.UpdatePerson(ctx, "p2", "New", nil, model.GenderFemale, &empty)
		assert.NoError(t, err)
		assert.Nil(t, p.PhotoPath)

		mockRepo.AssertExpectations(t)
		mockPhoto.AssertExpectations(t)
	})

	t.Run("Success - Replace Photo", func(t *testing.T) {
		oldPhoto := "old.jpg"
		p := &model.Person{ID: "p3", Name: "Old", PhotoPath: &oldPhoto}

		mockRepo.On("GetPersonByID", ctx, "p3").Return(p, nil).Once()
		mockPhoto.On("DeletePhoto", "old.jpg").Return(nil).Once()
		mockRepo.On("UpdatePerson", ctx, mock.AnythingOfType("*model.Person")).Return(nil).Once()

		newPhoto := "new.jpg"
		err := svc.UpdatePerson(ctx, "p3", "New", nil, model.GenderFemale, &newPhoto)
		assert.NoError(t, err)
		assert.Equal(t, "new.jpg", *p.PhotoPath)

		mockRepo.AssertExpectations(t)
		mockPhoto.AssertExpectations(t)
	})
	
	t.Run("Error - Person Not Found", func(t *testing.T) {
		mockRepo.On("GetPersonByID", ctx, "p4").Return(nil, errors.New("not found")).Once()
		
		err := svc.UpdatePerson(ctx, "p4", "New", nil, model.GenderFemale, nil)
		assert.Error(t, err)
		
		mockRepo.AssertExpectations(t)
	})
}

func TestPersonService_DeletePerson(t *testing.T) {
	mockRepo := new(MockPersonRepository)
	mockPhoto := new(MockPhotoStorage)
	svc := service.NewPersonService(mockRepo, mockPhoto)
	ctx := context.Background()

	t.Run("Success With Photo", func(t *testing.T) {
		photo := "pic.jpg"
		p := &model.Person{ID: "p1", PhotoPath: &photo}

		mockRepo.On("GetPersonByID", ctx, "p1").Return(p, nil).Once()
		mockRepo.On("DeletePerson", ctx, "p1").Return(nil).Once()
		mockPhoto.On("DeletePhoto", "pic.jpg").Return(nil).Once()

		err := svc.DeletePerson(ctx, "p1")
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockPhoto.AssertExpectations(t)
	})

	t.Run("Success Without Photo", func(t *testing.T) {
		p := &model.Person{ID: "p2", PhotoPath: nil}

		mockRepo.On("GetPersonByID", ctx, "p2").Return(p, nil).Once()
		mockRepo.On("DeletePerson", ctx, "p2").Return(nil).Once()

		err := svc.DeletePerson(ctx, "p2")
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockPhoto.AssertExpectations(t)
	})
}
