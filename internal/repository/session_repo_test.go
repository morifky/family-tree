package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"brayat/internal/model"
	"brayat/internal/repository"
)

func setupSessionDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Session{}, &model.AccessLink{})
	require.NoError(t, err)

	db.Exec("DELETE FROM access_links")
	db.Exec("DELETE FROM sessions")

	return db
}

func TestSessionRepository_CreateSession(t *testing.T) {
	db := setupSessionDB(t)
	repo := repository.NewSessionRepository(db)
	ctx := context.Background()

	session := &model.Session{
		Title:     "New Session",
		AdminCode: "admin123",
		ExpiresAt: time.Now().AddDate(0, 0, 30),
	}

	err := repo.CreateSession(ctx, session)
	assert.NoError(t, err)
	assert.NotEmpty(t, session.ID) // Assigned by gonanoid

	var s model.Session
	err = db.First(&s, "id = ?", session.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "New Session", s.Title)
}

func TestSessionRepository_GetSessionByID(t *testing.T) {
	db := setupSessionDB(t)
	repo := repository.NewSessionRepository(db)
	ctx := context.Background()

	session := &model.Session{ID: "sess1", Title: "t"}
	db.Create(session)

	s, err := repo.GetSessionByID(ctx, "sess1")
	assert.NoError(t, err)
	assert.Equal(t, "t", s.Title)

	_, err = repo.GetSessionByID(ctx, "fake")
	assert.Error(t, err)
}

func TestSessionRepository_GetSessionByAdminCode(t *testing.T) {
	db := setupSessionDB(t)
	repo := repository.NewSessionRepository(db)
	ctx := context.Background()

	session := &model.Session{ID: "sess2", Title: "t", AdminCode: "adm2"}
	db.Create(session)

	s, err := repo.GetSessionByAdminCode(ctx, "adm2")
	assert.NoError(t, err)
	assert.Equal(t, "sess2", s.ID)
}

func TestSessionRepository_UpdateSessionStatus(t *testing.T) {
	db := setupSessionDB(t)
	repo := repository.NewSessionRepository(db)
	ctx := context.Background()

	session := &model.Session{ID: "sess3", Title: "t", Status: model.SessionStatusActive}
	db.Create(session)

	err := repo.UpdateSessionStatus(ctx, "sess3", model.SessionStatusLocked)
	assert.NoError(t, err)

	var s model.Session
	db.First(&s, "id = ?", "sess3")
	assert.Equal(t, model.SessionStatusLocked, s.Status)
}

func TestSessionRepository_ExtendSessionExpiry(t *testing.T) {
	db := setupSessionDB(t)
	repo := repository.NewSessionRepository(db)
	ctx := context.Background()

	now := time.Now()
	session := &model.Session{ID: "sess4", Title: "t", ExpiresAt: now}
	db.Create(session)

	future := now.Add(24 * time.Hour)
	err := repo.ExtendSessionExpiry(ctx, "sess4", future)
	assert.NoError(t, err)

	var s model.Session
	db.First(&s, "id = ?", "sess4")
	assert.True(t, s.ExpiresAt.After(now))
}

func TestSessionRepository_AccessLinks(t *testing.T) {
	db := setupSessionDB(t)
	repo := repository.NewSessionRepository(db)
	ctx := context.Background()

	db.Create(&model.Session{ID: "sess5", Title: "test"})

	link := &model.AccessLink{
		SessionID: "sess5",
		Code:      "link123",
		Type:      model.AccessTypeEdit,
	}

	err := repo.CreateAccessLink(ctx, link)
	assert.NoError(t, err)
	assert.NotEmpty(t, link.ID)

	// Get by code
	l, err := repo.GetAccessLinkByCode(ctx, "link123")
	assert.NoError(t, err)
	assert.Equal(t, "sess5", l.SessionID)

	// Get array
	links, err := repo.GetAccessLinksBySessionID(ctx, "sess5")
	assert.NoError(t, err)
	assert.Len(t, links, 1)
}
