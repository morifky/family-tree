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
	ID        string           `gorm:"type:varchar(10);primaryKey" json:"id"`
	SessionID string           `gorm:"type:varchar(10);not null" json:"session_id"`
	Session   Session          `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE;" json:"-"`
	PersonAID string           `gorm:"type:varchar(10);not null" json:"person_a_id"`
	PersonA   Person           `gorm:"foreignKey:PersonAID;constraint:OnDelete:CASCADE;" json:"-"`
	PersonBID string           `gorm:"type:varchar(10);not null" json:"person_b_id"`
	PersonB   Person           `gorm:"foreignKey:PersonBID;constraint:OnDelete:CASCADE;" json:"-"`
	Type      RelationshipType `gorm:"type:varchar(20);not null" json:"type"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// BeforeCreate hook to auto-generate 10-character IDs
func (r *Relationship) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == "" {
		r.ID, err = gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 10)
	}
	return
}

// AfterSave touches the parent Session to update its UpdatedAt field
func (r *Relationship) AfterSave(tx *gorm.DB) (err error) {
	return tx.Table("sessions").Where("id = ?", r.SessionID).Update("updated_at", time.Now()).Error
}

// AfterDelete touches the parent Session to update its UpdatedAt field
func (r *Relationship) AfterDelete(tx *gorm.DB) (err error) {
	return tx.Table("sessions").Where("id = ?", r.SessionID).Update("updated_at", time.Now()).Error
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
