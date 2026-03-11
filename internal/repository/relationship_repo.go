package repository

import (
	"context"

	"brayat/internal/model"

	"gorm.io/gorm"
)

type relationshipRepository struct {
	db *gorm.DB
}

// NewRelationshipRepository creates a new SQLite-backed relationship repository.
func NewRelationshipRepository(db *gorm.DB) model.RelationshipRepository {
	return &relationshipRepository{db: db}
}

func (r *relationshipRepository) CreateRelationship(ctx context.Context, rel *model.Relationship) error {
	return r.db.WithContext(ctx).Create(rel).Error
}

func (r *relationshipRepository) GetRelationshipsBySessionID(ctx context.Context, sessionID string) ([]model.Relationship, error) {
	var relationships []model.Relationship
	if err := r.db.WithContext(ctx).Where("session_id = ?", sessionID).Find(&relationships).Error; err != nil {
		return nil, err
	}
	return relationships, nil
}

func (r *relationshipRepository) DeleteRelationship(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Relationship{}).Error
}

func (r *relationshipRepository) DeleteRelationshipsByPersonID(ctx context.Context, personID string) error {
	return r.db.WithContext(ctx).Where("person_a_id = ? OR person_b_id = ?", personID, personID).Delete(&model.Relationship{}).Error
}
