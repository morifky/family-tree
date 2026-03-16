package model

import (
	"context"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

// SessionStatus defines the allowed session statuses
type SessionStatus string

const (
	SessionStatusActive SessionStatus = "active"
	SessionStatusLocked SessionStatus = "locked"
	SessionStatusClosed SessionStatus = "closed"
)

// AccessType defines the allowed access link types
type AccessType string

const (
	AccessTypeAdmin AccessType = "admin"
	AccessTypeEdit  AccessType = "edit"
	AccessTypeView  AccessType = "view"
)

// Session represents a family tree session
type Session struct {
	ID        string        `gorm:"type:varchar(10);primaryKey" json:"id"`
	Title     string        `gorm:"type:varchar(100);not null" json:"title"`
	AdminCode string        `gorm:"type:varchar(10);uniqueIndex;not null" json:"admin_code"`
	Status    SessionStatus `gorm:"type:varchar(20);default:'active'" json:"status"`
	ExpiresAt time.Time     `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// AccessLink represents a shareable link that grants either view or edit access
type AccessLink struct {
	ID        string     `gorm:"type:varchar(10);primaryKey" json:"id"`
	SessionID string     `gorm:"type:varchar(10);not null" json:"session_id"`
	Session   Session    `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE;" json:"-"`
	Code      string     `gorm:"type:varchar(10);uniqueIndex;not null" json:"code"`
	Type      AccessType `gorm:"type:varchar(20);not null" json:"type"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// BeforeCreate hooks to auto-generate 10-character IDs
func (s *Session) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID, err = gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 10)
	}
	return
}

func (al *AccessLink) BeforeCreate(tx *gorm.DB) (err error) {
	if al.ID == "" {
		al.ID, err = gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 10)
	}
	return
}

// SessionRepository defines the expected behavior of a session data access layer
type SessionRepository interface {
	CreateSession(ctx context.Context, session *Session) error
	GetSessionByID(ctx context.Context, id string) (*Session, error)
	GetSessionByAdminCode(ctx context.Context, code string) (*Session, error)
	UpdateSessionStatus(ctx context.Context, id string, status SessionStatus) error
	ExtendSessionExpiry(ctx context.Context, id string, newExpiry time.Time) error

	CreateAccessLink(ctx context.Context, link *AccessLink) error
	GetAccessLinkByCode(ctx context.Context, code string) (*AccessLink, error)
	GetAccessLinksBySessionID(ctx context.Context, sessionID string) ([]AccessLink, error)
}

// SessionService defines the expected behavior of the business logic layer for sessions
type SessionService interface {
	CreateSession(ctx context.Context, title string) (*Session, error)
	GetSessionByID(ctx context.Context, id string) (*Session, error)
	GetSessionByAdminCode(ctx context.Context, code string) (*Session, error)
	UpdateSessionStatus(ctx context.Context, id string, status SessionStatus) error
	ExtendSessionExpiry(ctx context.Context, id string) error

	CreateAccessLink(ctx context.Context, sessionID string, accessType AccessType) (*AccessLink, error)
	VerifyAccessCode(ctx context.Context, code string) (sessionID string, access AccessType, err error)
	GetAccessLinksBySessionID(ctx context.Context, sessionID string) ([]AccessLink, error)
}
