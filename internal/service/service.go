package service

import (
	"brayat/internal/model"
	"brayat/internal/repository"
)

// Services contains all the service interfaces handling business logic.
type Services struct {
	Session      model.SessionService
	Person       model.PersonService
	Relationship model.RelationshipService
}

// NewServices initializes the service layer merging business logic handlers globally.
func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Session: NewSessionService(repos.Session),
		// We'll populate Person and Relationship services in MOR-14 and MOR-15
	}
}
