package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"brayat/internal/model"
	"brayat/internal/service"
)

// MockRelationshipRepository is a mock of RelationshipRepository
type MockRelationshipRepository struct {
	mock.Mock
}

func (m *MockRelationshipRepository) CreateRelationship(ctx context.Context, rel *model.Relationship) error {
	args := m.Called(ctx, rel)
	return args.Error(0)
}

func (m *MockRelationshipRepository) GetRelationshipsBySessionID(ctx context.Context, sessionID string) ([]model.Relationship, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) != nil {
		return args.Get(0).([]model.Relationship), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRelationshipRepository) DeleteRelationship(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRelationshipRepository) DeleteRelationshipsByPersonID(ctx context.Context, personID string) error {
	args := m.Called(ctx, personID)
	return args.Error(0)
}

func TestRelationshipService_CreateRelationship(t *testing.T) {
	mockRepo := new(MockRelationshipRepository)
	svc := service.NewRelationshipService(mockRepo)

	ctx := context.Background()
	sessionID := "sess123"
	personA := "personA"
	personB := "personB"
	relType := model.RelationshipTypeParentChild

	// Success case
	mockRepo.On("CreateRelationship", ctx, mock.AnythingOfType("*model.Relationship")).Return(nil).Once()

	rel, err := svc.CreateRelationship(ctx, sessionID, personA, personB, relType)

	assert.NoError(t, err)
	assert.NotNil(t, rel)
	assert.Equal(t, sessionID, rel.SessionID)
	assert.Equal(t, personA, rel.PersonAID)
	assert.Equal(t, personB, rel.PersonBID)
	assert.Equal(t, relType, rel.Type)

	// Error case
	expectedErr := errors.New("db error")
	mockRepo.On("CreateRelationship", ctx, mock.AnythingOfType("*model.Relationship")).Return(expectedErr).Once()

	rel2, err2 := svc.CreateRelationship(ctx, sessionID, personA, personB, relType)

	assert.ErrorIs(t, err2, expectedErr)
	assert.Nil(t, rel2)

	mockRepo.AssertExpectations(t)
}

func TestRelationshipService_GetRelationshipsBySessionID(t *testing.T) {
	mockRepo := new(MockRelationshipRepository)
	svc := service.NewRelationshipService(mockRepo)

	ctx := context.Background()
	sessionID := "sess123"
	expectedRels := []model.Relationship{
		{ID: "rel1", SessionID: sessionID},
	}

	// Success case
	mockRepo.On("GetRelationshipsBySessionID", ctx, sessionID).Return(expectedRels, nil).Once()

	rels, err := svc.GetRelationshipsBySessionID(ctx, sessionID)

	assert.NoError(t, err)
	assert.Equal(t, expectedRels, rels)

	// Error case
	expectedErr := errors.New("not found")
	mockRepo.On("GetRelationshipsBySessionID", ctx, sessionID).Return(nil, expectedErr).Once()

	rels2, err2 := svc.GetRelationshipsBySessionID(ctx, sessionID)

	assert.ErrorIs(t, err2, expectedErr)
	assert.Nil(t, rels2)

	mockRepo.AssertExpectations(t)
}

func TestRelationshipService_DeleteRelationship(t *testing.T) {
	mockRepo := new(MockRelationshipRepository)
	svc := service.NewRelationshipService(mockRepo)

	ctx := context.Background()
	relID := "rel123"

	// Success case
	mockRepo.On("DeleteRelationship", ctx, relID).Return(nil).Once()

	err := svc.DeleteRelationship(ctx, relID)

	assert.NoError(t, err)

	// Error case
	expectedErr := errors.New("delete error")
	mockRepo.On("DeleteRelationship", ctx, relID).Return(expectedErr).Once()

	err2 := svc.DeleteRelationship(ctx, relID)

	assert.ErrorIs(t, err2, expectedErr)

	mockRepo.AssertExpectations(t)
}
