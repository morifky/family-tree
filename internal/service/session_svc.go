package service

import (
	"context"
	"errors"
	"time"

	"brayat/internal/model"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type sessionService struct {
	repo model.SessionRepository
}

// NewSessionService returns a new SessionService object backed by a repository.
func NewSessionService(repo model.SessionRepository) model.SessionService {
	return &sessionService{repo: repo}
}

func (s *sessionService) CreateSession(ctx context.Context, title string) (*model.Session, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	code, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 10)
	if err != nil {
		return nil, err
	}

	session := &model.Session{
		Title:     title,
		AdminCode: code,
		ExpiresAt: time.Now().AddDate(0, 0, 30), // Default 30 days expiry
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *sessionService) GetSessionByID(ctx context.Context, id string) (*model.Session, error) {
	return s.repo.GetSessionByID(ctx, id)
}

func (s *sessionService) GetSessionByAdminCode(ctx context.Context, code string) (*model.Session, error) {
	return s.repo.GetSessionByAdminCode(ctx, code)
}

func (s *sessionService) UpdateSessionStatus(ctx context.Context, id string, status model.SessionStatus) error {
	switch status {
	case model.SessionStatusActive, model.SessionStatusLocked, model.SessionStatusClosed:
		// all valid
	default:
		return errors.New("invalid session status")
	}

	return s.repo.UpdateSessionStatus(ctx, id, status)
}

func (s *sessionService) ExtendSessionExpiry(ctx context.Context, id string) error {
	newExpiry := time.Now().AddDate(0, 0, 30) // Extend by another 30 days initially
	return s.repo.ExtendSessionExpiry(ctx, id, newExpiry)
}

func (s *sessionService) CreateAccessLink(ctx context.Context, sessionID string, accessType model.AccessType) (*model.AccessLink, error) {
	if accessType != model.AccessTypeEdit && accessType != model.AccessTypeView {
		return nil, errors.New("invalid access type")
	}

	code, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 10)
	if err != nil {
		return nil, err
	}

	link := &model.AccessLink{
		SessionID: sessionID,
		Code:      code,
		Type:      accessType,
	}

	if err := s.repo.CreateAccessLink(ctx, link); err != nil {
		return nil, err
	}

	return link, nil
}

func (s *sessionService) VerifyAccessCode(ctx context.Context, code string) (string, model.AccessType, error) {
	// First check if it's an admin code directly from the session
	session, err := s.repo.GetSessionByAdminCode(ctx, code)
	if err == nil {
		// Valid tracking code as an admin (which resolves to admin capability)
		return session.ID, model.AccessTypeAdmin, nil
	}

	// Maybe it's an access link code dynamically generated?
	link, err := s.repo.GetAccessLinkByCode(ctx, code)
	if err != nil {
		return "", "", errors.New("invalid code")
	}

	return link.SessionID, link.Type, nil
}

func (s *sessionService) GetAccessLinksBySessionID(ctx context.Context, sessionID string) ([]model.AccessLink, error) {
	return s.repo.GetAccessLinksBySessionID(ctx, sessionID)
}
