package repositories

import (
	"errors"
	"fmt"
	"sync"
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

type OrdersModel struct {
	ID        uint `gorm:"primaryKey"`
	OrderName string
	UserID    uint
}

type TestModelPreload struct {
	ID     uint          `gorm:"primaryKey"`
	Email  string        `gorm:"unique"`
	Orders []OrdersModel `gorm:"foreignKey:UserID"`
}

type TestModelWithStringID struct {
	ID    string `gorm:"primaryKey"`
	Email string `gorm:"unique"`
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

var mu sync.Mutex // Mutex to control concurrent access

func TestGenericRepository_Concurrent_Create(t *testing.T) {
	// Use file-based SQLite database to avoid in-memory issues
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	// Set SQLite busy timeout to handle concurrent locks
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get db from gorm: %v", err)
	}
	sqlDB.SetConnMaxLifetime(0) // Disable timeout
	sqlDB.SetMaxOpenConns(1)    // Ensure a single open connection to avoid concurrency issues
	sqlDB.SetMaxIdleConns(1)    // Only allow 1 idle connection

	// Ensure the migration happens
	err = db.AutoMigrate(&TestModel{})
	if err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewGenericRepository[TestModel, uint](db)

	var wg sync.WaitGroup
	concurrency := 10

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			mu.Lock() // Lock the mutex before accessing the database
			defer mu.Unlock()

			model := TestModel{Email: fmt.Sprintf("test%d@example.com", i)}
			_, err := repo.Create(model)
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	// Ensure that all records were created in the database
	var count int64
	err = db.Model(&TestModel{}).Count(&count).Error
	if err != nil {
		t.Fatalf("Error counting records: %v", err)
	}

	// Verify that the number of created records matches the concurrency level
	assert.Equal(t, int64(concurrency), count)
}

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

	model := TestModel{Email: "ToBeDeleted"}
	createdModel, err := repo.Create(model)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	err = repo.Delete(createdModel.ID, true)
	if err != nil {
		t.Fatalf("Failed to delete model permanently: %v", err)
	}

	var result TestModel
	err = db.Unscoped().First(&result, createdModel.ID).Error
	if err == nil {
		t.Errorf("Expected record to be deleted, but it still exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Expected 'record not found' error, but got: %v", err)
	}
}

func TestGenericRepository_Delete_DBError(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	model := TestModel{Email: "error@test.com"}
	_, err := repo.Create(model)
	if err != nil {
		t.Fatalf("Failed to set up test data: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.Close()

	err = repo.Delete(model.ID, true)
	if err == nil {
		t.Fatalf("Expected an error due to DB failure, but got none")
	}
}
func TestGenericRepository_Delete_NonExistentItem(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	nonExistentID := uint(9999)

	err := repo.Delete(nonExistentID, true)
	if err == nil {
		t.Fatalf("Expected an error when deleting a non-existent item, but got none")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("Expected 'record not found' error, but got: %v", err)
	}
}

func TestGenericRepository_Get_DBError(t *testing.T) {
	db := setupGORMSqlite(t, &TestModel{})
	repo := NewGenericRepository[TestModel, uint](db)

	// Create a model in the DB
	model := TestModel{Email: "error@test.com"}
	_, err := repo.Create(model)
	if err != nil {
		t.Fatalf("Failed to create test model: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.Close()

	_, err = repo.Get(model.ID, "")
	if err == nil {
		t.Fatalf("Expected an error due to DB failure, but got none")
	}
}

func TestGenericRepository_Get_WithPreload(t *testing.T) {
	db := setupGORMSqlite(t, &TestModelPreload{}, &OrdersModel{})
	repo := NewGenericRepository[TestModelPreload, uint](db)

	user := TestModelPreload{Email: "user@test.com"}
	err := db.Create(&user).Error
	if err != nil {
		t.Fatalf("Unexpected error creating user: %v", err)
	}

	order := OrdersModel{OrderName: "Test Order", UserID: user.ID}
	err = db.Create(&order).Error
	if err != nil {
		t.Fatalf("Unexpected error creating order: %v", err)
	}

	result, err := repo.Get(user.ID, "Orders")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result.Orders) == 0 {
		t.Fatalf("Expected preloaded orders, but got none")
	}

	if result.Orders[0].OrderName != "Test Order" {
		t.Fatalf("Expected order name 'Test Order', but got %v", result.Orders[0].OrderName)
	}
}

func TestGenericRepository_Delete_StringID(t *testing.T) {
	db := setupGORMSqlite(t, &TestModelWithStringID{})
	repo := NewGenericRepository[TestModelWithStringID, string](db)

	model := TestModelWithStringID{ID: "abc123", Email: "test@test.com"}
	err := db.Create(&model).Error
	if err != nil {
		t.Fatalf("Unexpected error creating model: %v", err)
	}

	err = repo.Delete(model.ID, false)
	if err != nil {
		t.Fatalf("Unexpected error during delete: %v", err)
	}

	var deletedModel TestModelWithStringID
	result := db.First(&deletedModel, "id = ?", model.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		t.Fatalf("Expected record to be deleted, but got unexpected error: %v", result.Error)
	}
}
