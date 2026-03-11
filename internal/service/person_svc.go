package service

import (
	"context"

	"brayat/internal/model"
	"brayat/internal/storage"
)

type personService struct {
	repo         model.PersonRepository
	photoStorage storage.PhotoStorage
}

// NewPersonService initializes a new PersonService with its dependencies.
func NewPersonService(repo model.PersonRepository, photoStorage storage.PhotoStorage) model.PersonService {
	return &personService{
		repo:         repo,
		photoStorage: photoStorage,
	}
}

func (s *personService) CreatePerson(ctx context.Context, sessionID string, name string, nickname *string, gender model.Gender, photoPath *string) (*model.Person, error) {
	person := &model.Person{
		SessionID: sessionID,
		Name:      name,
		Nickname:  nickname,
		Gender:    gender,
		PhotoPath: photoPath,
	}

	if err := s.repo.CreatePerson(ctx, person); err != nil {
		return nil, err
	}

	return person, nil
}

func (s *personService) GetPersonByID(ctx context.Context, id string) (*model.Person, error) {
	return s.repo.GetPersonByID(ctx, id)
}

func (s *personService) GetPeopleBySessionID(ctx context.Context, sessionID string) ([]model.Person, error) {
	return s.repo.GetPeopleBySessionID(ctx, sessionID)
}

func (s *personService) UpdatePerson(ctx context.Context, id string, name string, nickname *string, gender model.Gender, photoPath *string) error {
	person, err := s.repo.GetPersonByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete old photo if new photo is provided and is different
	if photoPath != nil && person.PhotoPath != nil && *photoPath != *person.PhotoPath {
		if s.photoStorage != nil {
			_ = s.photoStorage.DeletePhoto(*person.PhotoPath)
		}
	}

	// Update fields
	person.Name = name
	person.Nickname = nickname
	person.Gender = gender
	// Only update photoPath if explicitly requested (e.g., provided as not nil).
	// If the user wants to remove the photo entirely, they could pass a pointer to an empty string.
	if photoPath != nil {
		if *photoPath == "" {
			// They want to remove the photo. Delete current.
			if person.PhotoPath != nil && s.photoStorage != nil {
				_ = s.photoStorage.DeletePhoto(*person.PhotoPath)
			}
			person.PhotoPath = nil
		} else {
			person.PhotoPath = photoPath
		}
	}

	return s.repo.UpdatePerson(ctx, person)
}

func (s *personService) DeletePerson(ctx context.Context, id string) error {
	person, err := s.repo.GetPersonByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.repo.DeletePerson(ctx, id); err != nil {
		return err
	}

	// Delete photo if exists
	if person.PhotoPath != nil && s.photoStorage != nil {
		_ = s.photoStorage.DeletePhoto(*person.PhotoPath)
	}

	return nil
}
