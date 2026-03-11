package service

import (
	"brayat/internal/model"
	"brayat/internal/repository"
	"brayat/internal/storage"
)

// Services contains all the service interfaces handling business logic.
type Services struct {
	Session      model.SessionService
	Person       model.PersonService
	Relationship model.RelationshipService
}

// NewServices initializes the service layer merging business logic handlers globally.
func NewServices(repos *repository.Repositories, photoStorage storage.PhotoStorage) *Services {
	return &Services{
		Session:      NewSessionService(repos.Session),
		Person:       NewPersonService(repos.Person, photoStorage),
		Relationship: NewRelationshipService(repos.Relationship),
	}
}
