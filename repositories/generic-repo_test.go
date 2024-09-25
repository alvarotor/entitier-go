package repositories

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestModel struct {
	ID    uint   `gorm:"primaryKey"`
	Email string `gorm:"unique"`
}

type TestModelWithVariousFields struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique"`
	Age      int
	Salary   float64
	JoinedAt time.Time
}

func TestGenericRepository_Create_WithVariousFields(t *testing.T) {
	db := setupGORMSqlite(t, &TestModelWithVariousFields{})
	repo := NewGenericRepository[TestModelWithVariousFields, uint](db)

	model := TestModelWithVariousFields{
		Email:    "test@example.com",
		Age:      30,
		Salary:   50000.00,
		JoinedAt: time.Now(),
	}

	_, err := repo.Create(model)
	assert.NoError(t, err)
}

// setupGORMSqlite creates an in-memory SQLite database and returns a *gorm.DB connection
func setupGORMSqlite(t *testing.T, models ...interface{}) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// Migrate the schema for provided models
	for _, model := range models {
		err = db.AutoMigrate(model)
		if err != nil {
			t.Fatalf("failed to migrate database schema for model %v: %v", model, err)
		}
	}

	return db
}

func TestGenericRepository_Create(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	model := TestModel{Email: "Test"}

	// Act
	createdModel, err := repo.Create(model)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, model.Email, createdModel.Email)
}

func TestGenericRepository_GetAll(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	// Insert some records to test
	db.Create(&TestModel{Email: "Test1"})
	db.Create(&TestModel{Email: "Test2"})

	// Act
	result, err := repo.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
}

func TestGenericRepository_Get(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	// Insert a record to test
	db.Create(&TestModel{Email: "Test"})

	// Act
	result, err := repo.Get(1, "")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, uint(1), result.ID)
}

func TestGenericRepository_Update(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	// Insert a record to test
	db.Create(&TestModel{Email: "OldEmail"})

	updatedModel := TestModel{Email: "NewEmail"}

	// Act
	err := repo.Update(1, updatedModel)

	// Assert
	assert.NoError(t, err)

	var fetchedModel TestModel
	db.First(&fetchedModel, 1)
	assert.Equal(t, "NewEmail", fetchedModel.Email)
}

func TestGenericRepository_Delete(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	// Insert a record to test
	db.Create(&TestModel{Email: "ToDelete"})

	// Act
	err := repo.Delete(1, false)

	// Assert
	assert.NoError(t, err)

	var result TestModel
	tx := db.First(&result, 1)
	assert.Error(t, tx.Error) // Should return not found
}

func TestGenericRepository_Delete_Permanently(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	// Insert a record to test
	db.Create(&TestModel{Email: "ToDelete"})

	// Act
	err := repo.Delete(1, true)

	// Assert
	assert.NoError(t, err)

	var result TestModel
	tx := db.Unscoped().First(&result, 1)
	assert.Error(t, tx.Error) // Should return not found
}

func TestGenericRepository_Create_Error(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	model := TestModel{Email: "Test"}

	// Insert the same record twice to trigger a unique constraint violation
	db.Create(&model)

	// Attempt to create again and expect a duplicate error
	_, err := repo.Create(model)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UNIQUE constraint failed")
}

func TestGenericRepository_GetAll_Empty(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	// Act
	result, err := repo.GetAll()

	// Assert
	assert.Error(t, err)                          // Expecting an error
	assert.Equal(t, "no rows found", err.Error()) // Adjust based on actual error message
	assert.Equal(t, 0, len(result))               // Expecting empty result
}

func TestGenericRepository_Get_NotFound(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	// Try to get a record that doesn't exist
	_, err := repo.Get(999, "")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
}

func TestGenericRepository_Delete_NotFound(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	// Try to delete a record that doesn't exist
	err := repo.Delete(999, false)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
}

func TestGenericRepository_Delete_Permanent_NotFound(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	// Try to permanently delete a record that doesn't exist
	err := repo.Delete(999, true)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
}

func TestGenericRepository_Update_Error(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	nonExistentModel := TestModel{Email: "NonExistent"}

	// Try updating a non-existent record
	err := repo.Update(999, nonExistentModel)

	// Assert that an error is returned for trying to update a non-existent record
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
}

func TestGenericRepository_Create_NoRowsAffected(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	// Create and rollback within a transaction
	tx := db.Begin()
	defer tx.Rollback() // Ensure rollback happens even if the test fails

	// Create a model within the transaction
	model := TestModel{Email: "TransactionalModel"}
	tx.Create(&model)

	// Rollback the transaction so no changes are committed
	tx.Rollback() // No rows should be affected

	// Now create the model normally, should handle rollback
	createdModel, err := repo.Create(model)

	// Assert that Create handles the rollback case properly
	assert.NoError(t, err)
	assert.Equal(t, model, createdModel)
}

func TestGenericRepository_Create_InvalidData(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	invalidModel := TestModel{Email: ""}
	_, err := repo.Create(invalidModel)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "model cannot be empty")
}

func TestGenericRepository_Update_Partial(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	model := TestModel{Email: "InitialEmail"}
	createdModel, err := repo.Create(model)
	assert.NoError(t, err)

	updatedModel := TestModel{Email: "UpdatedEmail"}
	err = repo.Update(createdModel.ID, updatedModel)
	assert.NoError(t, err)

	fetchedModel, err := repo.Get(createdModel.ID, "")
	assert.NoError(t, err)
	assert.Equal(t, "UpdatedEmail", fetchedModel.Email)
}

func TestGenericRepository_SoftDelete(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	model := TestModel{Email: "ToDelete"}
	createdModel, err := repo.Create(model)
	assert.NoError(t, err)

	err = repo.Delete(createdModel.ID, false) // Soft delete
	assert.NoError(t, err)

	_, err = repo.Get(createdModel.ID, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
}

func TestGenericRepository_PermanentDelete_Success(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	// Create a model
	model := TestModel{Email: "ToBeDeleted"}
	createdModel, _ := repo.Create(model)

	// Permanent delete the model
	err := repo.Delete(createdModel.ID, true)
	assert.NoError(t, err)

	// Verify the record is deleted
	var result TestModel
	err = db.Unscoped().First(&result, createdModel.ID).Error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
}

// func TestGenericRepository_Concurrent_Create(t *testing.T) {
// 	db := setupGORMSqlite(t, &TestModel{})

// 	repo := NewGenericRepository[TestModel, uint](db)

// 	var wg sync.WaitGroup
// 	concurrency := 10

// 	for i := 0; i < concurrency; i++ {
// 		wg.Add(1)
// 		go func(i int) {
// 			defer wg.Done()
// 			model := TestModel{Email: fmt.Sprintf("test%d@example.com", i)}
// 			_, err := repo.Create(model)
// 			assert.NoError(t, err)
// 		}(i)
// 	}

// 	wg.Wait()
// }

func TestGenericRepository_GetAll_LargeDataSet(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	for i := 0; i < 1000; i++ {
		db.Create(&TestModel{Email: fmt.Sprintf("test%d@example.com", i)})
	}

	// Act
	result, err := repo.GetAll()

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1000, len(result))
}
func TestGenericRepository_Create_SpecialCharacters(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	model := TestModel{Email: "!@#$%^&*()_+-=~`"}
	_, err := repo.Create(model)
	assert.NoError(t, err)
}

func TestGenericRepository_Delete_Unscoped(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	// Create a record that will be deleted
	model := TestModel{Email: "ToBeDeleted"}
	createdModel, err := repo.Create(model)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	// Now attempt to delete it permanently
	err = repo.Delete(createdModel.ID, true)
	if err != nil {
		t.Fatalf("Failed to delete model permanently: %v", err)
	}

	// Verify that the record is deleted
	var result TestModel
	err = db.Unscoped().First(&result, createdModel.ID).Error
	if err == nil {
		t.Errorf("Expected record to be deleted, but it still exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Expected 'record not found' error, but got: %v", err)
	}
}
