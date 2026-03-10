package repository

import (
	"context"
	"time"

	"brayat/internal/model"

	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository creates a new SQLite-backed session repository
func NewSessionRepository(db *gorm.DB) model.SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) CreateSession(ctx context.Context, session *model.Session) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *sessionRepository) GetSessionByID(ctx context.Context, id string) (*model.Session, error) {
	var session model.Session
	err := r.db.WithContext(ctx).First(&session, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) GetSessionByAdminCode(ctx context.Context, code string) (*model.Session, error) {
	var session model.Session
	err := r.db.WithContext(ctx).First(&session, "admin_code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepository) UpdateSessionStatus(ctx context.Context, id string, status model.SessionStatus) error {
	return r.db.WithContext(ctx).Model(&model.Session{}).Where("id = ?", id).Update("status", status).Error
}

func (r *sessionRepository) ExtendSessionExpiry(ctx context.Context, id string, newExpiry time.Time) error {
	return r.db.WithContext(ctx).Model(&model.Session{}).Where("id = ?", id).Update("expires_at", newExpiry).Error
}

func (r *sessionRepository) CreateAccessLink(ctx context.Context, link *model.AccessLink) error {
	return r.db.WithContext(ctx).Create(link).Error
}

func (r *sessionRepository) GetAccessLinkByCode(ctx context.Context, code string) (*model.AccessLink, error) {
	var link model.AccessLink
	err := r.db.WithContext(ctx).Preload("Session").First(&link, "code = ?", code).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *sessionRepository) GetAccessLinksBySessionID(ctx context.Context, sessionID string) ([]model.AccessLink, error) {
	var links []model.AccessLink
	err := r.db.WithContext(ctx).Where("session_id = ?", sessionID).Find(&links).Error
	return links, err
}
