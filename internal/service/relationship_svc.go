package service

import (
	"context"

	"brayat/internal/model"
)

type relationshipService struct {
	repo model.RelationshipRepository
}

// NewRelationshipService initializes a new RelationshipService with its dependencies.
func NewRelationshipService(repo model.RelationshipRepository) model.RelationshipService {
	return &relationshipService{
		repo: repo,
	}
}

func (s *relationshipService) CreateRelationship(ctx context.Context, sessionID string, personAID string, personBID string, relType model.RelationshipType) (*model.Relationship, error) {
	rel := &model.Relationship{
		SessionID: sessionID,
		PersonAID: personAID,
		PersonBID: personBID,
		Type:      relType,
	}

	if err := s.repo.CreateRelationship(ctx, rel); err != nil {
		return nil, err
	}

	return rel, nil
}

func (s *relationshipService) GetRelationshipsBySessionID(ctx context.Context, sessionID string) ([]model.Relationship, error) {
	return s.repo.GetRelationshipsBySessionID(ctx, sessionID)
}

func (s *relationshipService) DeleteRelationship(ctx context.Context, id string) error {
	return s.repo.DeleteRelationship(ctx, id)
}
