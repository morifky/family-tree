package model

import (
	"context"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

// RelationshipType defines the allowed types of relationships
type RelationshipType string

const (
	RelationshipTypeParentChild RelationshipType = "parent_child" // PersonA is parent of PersonB
	RelationshipTypeSpouse      RelationshipType = "spouse"       // PersonA and PersonB are spouses/partners
)

// Relationship connects two people within a session
type Relationship struct {
	ID        string           `gorm:"type:varchar(10);primaryKey"`
	SessionID string           `gorm:"type:varchar(10);not null"`
	Session   Session          `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE;"`
	PersonAID string           `gorm:"type:varchar(10);not null"`
	PersonA   Person           `gorm:"foreignKey:PersonAID;constraint:OnDelete:CASCADE;"`
	PersonBID string           `gorm:"type:varchar(10);not null"`
	PersonB   Person           `gorm:"foreignKey:PersonBID;constraint:OnDelete:CASCADE;"`
	Type      RelationshipType `gorm:"type:varchar(20);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// BeforeCreate hook to auto-generate 10-character IDs
func (r *Relationship) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == "" {
		r.ID, err = gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 10)
	}
	return
}

// RelationshipRepository defines the expected behavior of a relationship data access layer
type RelationshipRepository interface {
	CreateRelationship(ctx context.Context, rel *Relationship) error
	GetRelationshipsBySessionID(ctx context.Context, sessionID string) ([]Relationship, error)
	DeleteRelationship(ctx context.Context, id string) error
	DeleteRelationshipsByPersonID(ctx context.Context, personID string) error
}

// RelationshipService defines the expected behavior of the business logic layer for relationships
type RelationshipService interface {
	CreateRelationship(ctx context.Context, sessionID string, personAID string, personBID string, relType RelationshipType) (*Relationship, error)
	GetRelationshipsBySessionID(ctx context.Context, sessionID string) ([]Relationship, error)
	DeleteRelationship(ctx context.Context, id string) error
}
