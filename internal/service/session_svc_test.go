package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"brayat/internal/model"
	"brayat/internal/service"
)

// MockSessionRepository is a mock of SessionRepository
type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) CreateSession(ctx context.Context, session *model.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepository) GetSessionByID(ctx context.Context, id string) (*model.Session, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Session), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSessionRepository) GetSessionByAdminCode(ctx context.Context, code string) (*model.Session, error) {
	args := m.Called(ctx, code)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Session), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSessionRepository) UpdateSessionStatus(ctx context.Context, id string, status model.SessionStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockSessionRepository) ExtendSessionExpiry(ctx context.Context, id string, newExpiry time.Time) error {
	args := m.Called(ctx, id, newExpiry)
	return args.Error(0)
}

func (m *MockSessionRepository) CreateAccessLink(ctx context.Context, link *model.AccessLink) error {
	args := m.Called(ctx, link)
	return args.Error(0)
}

func (m *MockSessionRepository) GetAccessLinkByCode(ctx context.Context, code string) (*model.AccessLink, error) {
	args := m.Called(ctx, code)
	if args.Get(0) != nil {
		return args.Get(0).(*model.AccessLink), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockSessionRepository) GetAccessLinksBySessionID(ctx context.Context, sessionID string) ([]model.AccessLink, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) != nil {
		return args.Get(0).([]model.AccessLink), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestSessionService_CreateSession(t *testing.T) {
	mockRepo := new(MockSessionRepository)
	svc := service.NewSessionService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("CreateSession", ctx, mock.AnythingOfType("*model.Session")).Return(nil).Once()

		s, err := svc.CreateSession(ctx, "Family Tree")
		assert.NoError(t, err)
		assert.Equal(t, "Family Tree", s.Title)
		assert.NotEmpty(t, s.AdminCode)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Empty Title", func(t *testing.T) {
		s, err := svc.CreateSession(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, s)
	})

	t.Run("DB Error", func(t *testing.T) {
		mockRepo.On("CreateSession", ctx, mock.AnythingOfType("*model.Session")).Return(errors.New("db error")).Once()
		
		s, err := svc.CreateSession(ctx, "Family Tree")
		assert.Error(t, err)
		assert.Nil(t, s)
		
		mockRepo.AssertExpectations(t)
	})
}

func TestSessionService_GetSession(t *testing.T) {
	mockRepo := new(MockSessionRepository)
	svc := service.NewSessionService(mockRepo)
	ctx := context.Background()

	s := &model.Session{ID: "sess1", Title: "t"}
	mockRepo.On("GetSessionByID", ctx, "sess1").Return(s, nil).Once()
	res1, err := svc.GetSessionByID(ctx, "sess1")
	assert.NoError(t, err)
	assert.Equal(t, "sess1", res1.ID)

	mockRepo.On("GetSessionByAdminCode", ctx, "adm1").Return(s, nil).Once()
	res2, err := svc.GetSessionByAdminCode(ctx, "adm1")
	assert.NoError(t, err)
	assert.Equal(t, "sess1", res2.ID)
}

func TestSessionService_UpdateSessionStatus(t *testing.T) {
	mockRepo := new(MockSessionRepository)
	svc := service.NewSessionService(mockRepo)
	ctx := context.Background()

	t.Run("Valid Status", func(t *testing.T) {
		mockRepo.On("UpdateSessionStatus", ctx, "sess1", model.SessionStatusLocked).Return(nil).Once()
		err := svc.UpdateSessionStatus(ctx, "sess1", model.SessionStatusLocked)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid Status", func(t *testing.T) {
		err := svc.UpdateSessionStatus(ctx, "sess1", "invalid_status")
		assert.Error(t, err)
	})
}

func TestSessionService_ExtendSessionExpiry(t *testing.T) {
	mockRepo := new(MockSessionRepository)
	svc := service.NewSessionService(mockRepo)
	ctx := context.Background()

	mockRepo.On("ExtendSessionExpiry", ctx, "sess1", mock.AnythingOfType("time.Time")).Return(nil).Once()
	err := svc.ExtendSessionExpiry(ctx, "sess1")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestSessionService_CreateAccessLink(t *testing.T) {
	mockRepo := new(MockSessionRepository)
	svc := service.NewSessionService(mockRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("CreateAccessLink", ctx, mock.AnythingOfType("*model.AccessLink")).Return(nil).Once()

		l, err := svc.CreateAccessLink(ctx, "sess1", model.AccessTypeEdit)
		assert.NoError(t, err)
		assert.Equal(t, model.AccessTypeEdit, l.Type)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid Type", func(t *testing.T) {
		l, err := svc.CreateAccessLink(ctx, "sess1", "wrong_type")
		assert.Error(t, err)
		assert.Nil(t, l)
	})

	t.Run("DB Error", func(t *testing.T) {
		mockRepo.On("CreateAccessLink", ctx, mock.AnythingOfType("*model.AccessLink")).Return(errors.New("db error")).Once()
		
		l, err := svc.CreateAccessLink(ctx, "sess1", model.AccessTypeEdit)
		assert.Error(t, err)
		assert.Nil(t, l)
		
		mockRepo.AssertExpectations(t)
	})
}

func TestSessionService_VerifyAccessCode(t *testing.T) {
	mockRepo := new(MockSessionRepository)
	svc := service.NewSessionService(mockRepo)
	ctx := context.Background()

	t.Run("Admin Code", func(t *testing.T) {
		s := &model.Session{ID: "sess1", Title: "t"}
		mockRepo.On("GetSessionByAdminCode", ctx, "adm1").Return(s, nil).Once()

		sessionID, accessType, err := svc.VerifyAccessCode(ctx, "adm1")
		assert.NoError(t, err)
		assert.Equal(t, "sess1", sessionID)
		assert.Equal(t, model.AccessTypeAdmin, accessType)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Access Link Code", func(t *testing.T) {
		mockRepo.On("GetSessionByAdminCode", ctx, "link1").Return(nil, errors.New("not admin code")).Once()
		
		l := &model.AccessLink{SessionID: "sess2", Type: model.AccessTypeView}
		mockRepo.On("GetAccessLinkByCode", ctx, "link1").Return(l, nil).Once()

		sessionID, accessType, err := svc.VerifyAccessCode(ctx, "link1")
		assert.NoError(t, err)
		assert.Equal(t, "sess2", sessionID)
		assert.Equal(t, model.AccessTypeView, accessType)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Invalid Code", func(t *testing.T) {
		mockRepo.On("GetSessionByAdminCode", ctx, "wrong").Return(nil, errors.New("not admin code")).Once()
		mockRepo.On("GetAccessLinkByCode", ctx, "wrong").Return(nil, errors.New("not link code")).Once()

		id, typ, err := svc.VerifyAccessCode(ctx, "wrong")
		assert.Error(t, err)
		assert.Empty(t, id)
		assert.Empty(t, typ)

		mockRepo.AssertExpectations(t)
	})
}

func TestSessionService_GetAccessLinksBySessionID(t *testing.T) {
	mockRepo := new(MockSessionRepository)
	svc := service.NewSessionService(mockRepo)
	ctx := context.Background()

	expectedLinks := []model.AccessLink{{ID: "link1"}}
	mockRepo.On("GetAccessLinksBySessionID", ctx, "sess1").Return(expectedLinks, nil).Once()

	links, err := svc.GetAccessLinksBySessionID(ctx, "sess1")
	assert.NoError(t, err)
	assert.Len(t, links, 1)

	mockRepo.AssertExpectations(t)
}
