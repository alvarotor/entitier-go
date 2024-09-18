package repositories

import (
	"errors"
	"log"
	"testing"

	"github.com/alvarotor/entitier-go/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MockModel struct {
	ID    uint   `gorm:"primaryKey"`
	Email string `gorm:"unique"`
	Name  string
}

// Setup the in-memory SQLite database for testing purposes.
func setupDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Automatically migrate the mock model
	err = db.AutoMigrate(&MockModel{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestGenericRepository_Create(t *testing.T) {
	db, err := setupDB()
	assert.NoError(t, err)

	repo := NewGenericRepository[MockModel, uint](db)

	// Test successful creation
	model := MockModel{Name: "John Doe", Email: "john@example.com"}
	created, err := repo.Create(model)
	assert.NoError(t, err)
	assert.Equal(t, model.Name, created.Name)

	// Test duplicate key error (violating unique constraint on Email)
	_, err = repo.Create(model)
	assert.Error(t, err)
	log.Println("err.Error()")
	log.Println(err.Error())
	log.Println(models.ErrDuplicatedKeyEmail)
	log.Println(gorm.ErrDuplicatedKey)
	assert.Equal(t, err.Error(), "UNIQUE constraint failed: mock_models.email")
	// assert.True(t, errors.Is(err, models.ErrDuplicatedKeyEmail))
	// assert.True(t, errors.Is(err, gorm.ErrDuplicatedKey))
}

func TestGenericRepository_GetAll(t *testing.T) {
	db, err := setupDB()
	assert.NoError(t, err)

	repo := NewGenericRepository[MockModel, uint](db)

	// Initially should return models.ErrNotFound
	all, err := repo.GetAll()
	assert.Error(t, err)
	assert.True(t, errors.Is(err, models.ErrNotFound))

	// Create some mock records
	model1 := MockModel{Name: "John Doe", Email: "john@example.com"}
	model2 := MockModel{Name: "Jane Smith", Email: "jane@example.com"}
	_, err = repo.Create(model1)
	assert.NoError(t, err)
	_, err = repo.Create(model2)
	assert.NoError(t, err)

	// Fetch all records
	all, err = repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, all, 2)
}

func TestGenericRepository_Get(t *testing.T) {
	db, err := setupDB()
	assert.NoError(t, err)

	repo := NewGenericRepository[MockModel, uint](db)

	// Test getting a non-existent record
	_, err = repo.Get(1, "")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

	// Create a mock record
	model := MockModel{Name: "John Doe", Email: "john@example.com"}
	created, err := repo.Create(model)
	assert.NoError(t, err)

	// Test fetching the created record
	fetched, err := repo.Get(created.ID, "")
	assert.NoError(t, err)
	assert.Equal(t, created.Name, fetched.Name)
}

func TestGenericRepository_Update(t *testing.T) {
	db, err := setupDB()
	assert.NoError(t, err)

	repo := NewGenericRepository[MockModel, uint](db)

	// Test updating a non-existent record
	nonExistentModel := MockModel{ID: 999, Name: "Doesn't Exist", Email: "noone@example.com"}
	err = repo.Update(nonExistentModel.ID, nonExistentModel)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

	// Create a mock record
	model := MockModel{Name: "John Doe", Email: "john@example.com"}
	created, err := repo.Create(model)
	assert.NoError(t, err)

	// Test successful update
	created.Name = "John Updated"
	err = repo.Update(created.ID, created)
	assert.NoError(t, err)

	// Verify the update
	updated, err := repo.Get(created.ID, "")
	assert.NoError(t, err)
	assert.Equal(t, "John Updated", updated.Name)
}

func TestGenericRepository_Delete(t *testing.T) {
	db, err := setupDB()
	assert.NoError(t, err)

	repo := NewGenericRepository[MockModel, uint](db)

	// Test deleting a non-existent record
	err = repo.Delete(999, false)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

	// Create a mock record
	model := MockModel{Name: "John Doe", Email: "john@example.com"}
	created, err := repo.Create(model)
	assert.NoError(t, err)

	// Test soft delete
	err = repo.Delete(created.ID, false)
	assert.NoError(t, err)

	// Verify that the record is soft deleted (not found)
	_, err = repo.Get(created.ID, "")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

	// Test hard delete
	created, err = repo.Create(MockModel{Name: "Jane Doe", Email: "jane@example.com"})
	assert.NoError(t, err)

	err = repo.Delete(created.ID, true)
	assert.NoError(t, err)

	// Verify that the record is permanently deleted
	_, err = repo.Get(created.ID, "")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
