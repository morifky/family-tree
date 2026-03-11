package repository

import (
	"brayat/internal/model"

	"gorm.io/gorm"
)

// Repositories contains all the repository interfaces holding business data storage logic.
type Repositories struct {
	Session      model.SessionRepository
	Person       model.PersonRepository
	Relationship model.RelationshipRepository
}

// NewRepositories initializes and returns a merged Repositories structure linking to SQLite.
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Session: NewSessionRepository(db),
		Person:  NewPersonRepository(db),
		// We'll populate Relationship repo in MOR-15.
	}
}
