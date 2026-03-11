package repository

import (
	"context"

	"brayat/internal/model"

	"gorm.io/gorm"
)

type personRepository struct {
	db *gorm.DB
}

// NewPersonRepository creates a new SQLite-backed person repository.
func NewPersonRepository(db *gorm.DB) model.PersonRepository {
	return &personRepository{db: db}
}

func (r *personRepository) CreatePerson(ctx context.Context, person *model.Person) error {
	return r.db.WithContext(ctx).Create(person).Error
}

func (r *personRepository) GetPersonByID(ctx context.Context, id string) (*model.Person, error) {
	var person model.Person
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&person).Error; err != nil {
		return nil, err
	}
	return &person, nil
}

func (r *personRepository) GetPeopleBySessionID(ctx context.Context, sessionID string) ([]model.Person, error) {
	var people []model.Person
	if err := r.db.WithContext(ctx).Where("session_id = ?", sessionID).Order("created_at ASC").Find(&people).Error; err != nil {
		return nil, err
	}
	return people, nil
}

func (r *personRepository) UpdatePerson(ctx context.Context, person *model.Person) error {
	// gorm will update all non-zero fields, or we can force all fields by passing the struct
	return r.db.WithContext(ctx).Save(person).Error
}

func (r *personRepository) DeletePerson(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Person{}).Error
}
