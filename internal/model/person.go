package model

import (
	"context"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

// Gender defines the allowed gender options
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// Person represents an individual in the family tree
type Person struct {
	ID        string  `gorm:"type:varchar(10);primaryKey"`
	SessionID string  `gorm:"type:varchar(10);not null"`
	Session   Session `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE;"`
	Name      string  `gorm:"type:varchar(100);not null"`
	Nickname  *string `gorm:"type:varchar(50)"`
	Gender    Gender  `gorm:"type:varchar(10);not null"`
	PhotoPath *string `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate hook to auto-generate 10-character IDs
func (p *Person) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID, err = gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 10)
	}
	return
}

// AfterSave touches the parent Session to update its UpdatedAt field
func (p *Person) AfterSave(tx *gorm.DB) (err error) {
	return tx.Model(&Session{ID: p.SessionID}).Update("updated_at", time.Now()).Error
}

// AfterDelete touches the parent Session to update its UpdatedAt field
func (p *Person) AfterDelete(tx *gorm.DB) (err error) {
	return tx.Model(&Session{ID: p.SessionID}).Update("updated_at", time.Now()).Error
}

// PersonRepository defines the expected behavior of a person data access layer
type PersonRepository interface {
	CreatePerson(ctx context.Context, person *Person) error
	GetPersonByID(ctx context.Context, id string) (*Person, error)
	GetPeopleBySessionID(ctx context.Context, sessionID string) ([]Person, error)
	UpdatePerson(ctx context.Context, person *Person) error
	DeletePerson(ctx context.Context, id string) error
}

// PersonService defines the expected behavior of the business logic layer for people
type PersonService interface {
	CreatePerson(ctx context.Context, sessionID string, name string, nickname *string, gender Gender, photoPath *string) (*Person, error)
	GetPersonByID(ctx context.Context, id string) (*Person, error)
	GetPeopleBySessionID(ctx context.Context, sessionID string) ([]Person, error)
	UpdatePerson(ctx context.Context, id string, name string, nickname *string, gender Gender, photoPath *string) error
	DeletePerson(ctx context.Context, id string) error
}
