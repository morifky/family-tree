package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"brayat/internal/model"
	"brayat/internal/repository"
)

func setupPersonDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.Session{}, &model.Person{})
	require.NoError(t, err)

	db.Exec("DELETE FROM people")
	db.Exec("DELETE FROM sessions")

	return db
}

func TestPersonRepository_CreatePerson(t *testing.T) {
	db := setupPersonDB(t)
	repo := repository.NewPersonRepository(db)
	ctx := context.Background()

	session := model.Session{ID: "sessP1", Title: "test"}
	db.Create(&session)

	person := &model.Person{
		SessionID: "sessP1",
		Name:      "John Doe",
		Gender:    model.GenderMale,
	}

	err := repo.CreatePerson(ctx, person)
	assert.NoError(t, err)
	assert.NotEmpty(t, person.ID)

	var p model.Person
	err = db.First(&p, "id = ?", person.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", p.Name)
}

func TestPersonRepository_GetPersonByID(t *testing.T) {
	db := setupPersonDB(t)
	repo := repository.NewPersonRepository(db)
	ctx := context.Background()

	db.Create(&model.Session{ID: "sessP2", Title: "t"})
	person := &model.Person{ID: "pA", SessionID: "sessP2", Name: "Jane", Gender: model.GenderFemale}
	db.Create(person)

	p, err := repo.GetPersonByID(ctx, "pA")
	assert.NoError(t, err)
	assert.Equal(t, "Jane", p.Name)

	_, err = repo.GetPersonByID(ctx, "nonexistent")
	assert.Error(t, err)
}

func TestPersonRepository_GetPeopleBySessionID(t *testing.T) {
	db := setupPersonDB(t)
	repo := repository.NewPersonRepository(db)
	ctx := context.Background()

	db.Create(&model.Session{ID: "sessP3", Title: "t"})
	db.Create(&model.Person{ID: "p1", SessionID: "sessP3", Name: "A", Gender: model.GenderMale})
	db.Create(&model.Person{ID: "p2", SessionID: "sessP3", Name: "B", Gender: model.GenderFemale})

	people, err := repo.GetPeopleBySessionID(ctx, "sessP3")
	assert.NoError(t, err)
	assert.Len(t, people, 2)
}

func TestPersonRepository_UpdatePerson(t *testing.T) {
	db := setupPersonDB(t)
	repo := repository.NewPersonRepository(db)
	ctx := context.Background()

	db.Create(&model.Session{ID: "sessP4", Title: "t"})
	person := &model.Person{ID: "updateP", SessionID: "sessP4", Name: "Old", Gender: model.GenderMale}
	db.Create(person)

	person.Name = "New"
	err := repo.UpdatePerson(ctx, person)
	assert.NoError(t, err)

	var p model.Person
	db.First(&p, "id = ?", "updateP")
	assert.Equal(t, "New", p.Name)
}

func TestPersonRepository_DeletePerson(t *testing.T) {
	db := setupPersonDB(t)
	repo := repository.NewPersonRepository(db)
	ctx := context.Background()

	db.Create(&model.Session{ID: "sessP5", Title: "t"})
	db.Create(&model.Person{ID: "delP", SessionID: "sessP5", Name: "Delete Me", Gender: model.GenderMale})

	err := repo.DeletePerson(ctx, "delP")
	assert.NoError(t, err)

	var count int64
	db.Model(&model.Person{}).Where("id = ?", "delP").Count(&count)
	assert.Equal(t, int64(0), count)

	// Test deleting non-existent person
	err = repo.DeletePerson(ctx, "nonexistent")
	assert.NoError(t, err) // Should return nil if not found
}
