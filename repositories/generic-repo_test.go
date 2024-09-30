package repositories

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/alvarotor/entitier-go/mocks"
	"github.com/alvarotor/entitier-go/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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

var ctx = context.Background()

func TestGenericRepository_Create_WithVariousFields(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &TestModelWithVariousFields{})
	repo := NewGenericRepository[TestModelWithVariousFields, uint](db)

	model := TestModelWithVariousFields{
		Email:    "test@example.com",
		Age:      30,
		Salary:   50000.00,
		JoinedAt: time.Now(),
	}

	_, err := repo.Create(ctx, model)
	assert.NoError(t, err)
}

func TestGenericRepository_Create(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	model := mocks.TestModel{Email: "Test"}

	createdModel, err := repo.Create(ctx, model)

	assert.NoError(t, err)
	assert.Equal(t, model.Email, createdModel.Email)
}

func TestGenericRepository_GetAll(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	db.Create(&mocks.TestModel{Email: "Test1"})
	db.Create(&mocks.TestModel{Email: "Test2"})

	result, err := repo.GetAll(ctx)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
}

func TestGenericRepository_Get(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	db.Create(&mocks.TestModel{Email: "Test"})

	result, err := repo.Get(ctx, 1, "")

	assert.NoError(t, err)
	assert.Equal(t, uint(1), result.ID)
}

func TestGenericRepository_Update(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	db.Create(&mocks.TestModel{Email: "OldEmail"})

	updatedModel := mocks.TestModel{Email: "NewEmail"}

	err := repo.Update(ctx, 1, updatedModel)

	assert.NoError(t, err)

	var fetchedModel mocks.TestModel
	db.First(&fetchedModel, 1)
	assert.Equal(t, "NewEmail", fetchedModel.Email)
}

func TestGenericRepository_Delete(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	db.Create(&mocks.TestModel{Email: "ToDelete"})

	err := repo.Delete(ctx, 1, false)

	assert.NoError(t, err)

	var result mocks.TestModel
	tx := db.First(&result, 1)
	assert.Error(t, tx.Error)
}

func TestGenericRepository_Delete_Permanently(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	db.Create(&mocks.TestModel{Email: "ToDelete"})

	err := repo.Delete(ctx, 1, true)

	assert.NoError(t, err)

	var result mocks.TestModel
	tx := db.Unscoped().First(&result, 1)
	assert.Error(t, tx.Error)
}

func TestGenericRepository_Create_Error(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	model := mocks.TestModel{Email: "Test"}

	db.Create(&model)

	_, err := repo.Create(ctx, model)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UNIQUE constraint failed")
}

func TestGenericRepository_GetAll_Empty(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	result, err := repo.GetAll(ctx)

	assert.Error(t, err)
	assert.Equal(t, models.ErrNotFound, err)
	assert.Equal(t, 0, len(result))
}

func TestGenericRepository_Get_NotFound(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	_, err := repo.Get(ctx, 999, "")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
}

func TestGenericRepository_Delete_NotFound(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	err := repo.Delete(ctx, 999, false)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
}

func TestGenericRepository_Delete_Permanent_NotFound(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	err := repo.Delete(ctx, 999, true)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
}

func TestGenericRepository_Update_Error(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	nonExistentModel := mocks.TestModel{Email: "NonExistent"}

	err := repo.Update(ctx, 999, nonExistentModel)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
}

func TestGenericRepository_Create_NoRowsAffected(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	tx := db.Begin()
	defer tx.Rollback()

	model := mocks.TestModel{Email: "TransactionalModel"}
	tx.Create(&model)

	tx.Rollback()

	createdModel, err := repo.Create(ctx, model)

	assert.NoError(t, err)
	assert.Equal(t, model, createdModel)
}

func TestGenericRepository_Create_InvalidData(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	invalidModel := mocks.TestModel{Email: ""}
	_, err := repo.Create(ctx, invalidModel)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), models.ErrModelCannotBeEmpty.Error())
}

func TestGenericRepository_Update_Partial(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	model := mocks.TestModel{Email: "InitialEmail"}
	createdModel, err := repo.Create(ctx, model)
	assert.NoError(t, err)

	updatedModel := mocks.TestModel{Email: "UpdatedEmail"}
	err = repo.Update(ctx, createdModel.ID, updatedModel)
	assert.NoError(t, err)

	fetchedModel, err := repo.Get(ctx, createdModel.ID, "")
	assert.NoError(t, err)
	assert.Equal(t, "UpdatedEmail", fetchedModel.Email)
}

func TestGenericRepository_SoftDelete(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	model := mocks.TestModel{Email: "ToDelete"}
	createdModel, err := repo.Create(ctx, model)
	assert.NoError(t, err)

	err = repo.Delete(ctx, createdModel.ID, false)
	assert.NoError(t, err)

	_, err = repo.Get(ctx, createdModel.ID, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
}

func TestGenericRepository_PermanentDelete_Success(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	model := mocks.TestModel{Email: "ToBeDeleted"}
	createdModel, _ := repo.Create(ctx, model)

	err := repo.Delete(ctx, createdModel.ID, true)
	assert.NoError(t, err)

	var result mocks.TestModel
	err = db.Unscoped().First(&result, createdModel.ID).Error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "record not found")
}

var mu sync.Mutex

func TestGenericRepository_Concurrent_Create(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get db from gorm: %v", err)
	}
	sqlDB.SetConnMaxLifetime(0) // Disable timeout
	sqlDB.SetMaxOpenConns(1)    // Ensure a single open connection to avoid concurrency issues
	sqlDB.SetMaxIdleConns(1)    // Only allow 1 idle connection

	err = db.AutoMigrate(&mocks.TestModel{})
	if err != nil {
		t.Fatalf("failed to migrate database schema: %v", err)
	}

	repo := NewGenericRepository[mocks.TestModel, uint](db)

	var wg sync.WaitGroup
	concurrency := 10

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			mu.Lock()
			defer mu.Unlock()

			model := mocks.TestModel{Email: fmt.Sprintf("test%d@example.com", i)}
			_, err := repo.Create(ctx, model)
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	var count int64
	err = db.Model(&mocks.TestModel{}).Count(&count).Error
	if err != nil {
		t.Fatalf("Error counting records: %v", err)
	}

	assert.Equal(t, int64(concurrency), count)
}

func TestGenericRepository_GetAll_LargeDataSet(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	for i := 0; i < 1000; i++ {
		db.Create(&mocks.TestModel{Email: fmt.Sprintf("test%d@example.com", i)})
	}

	result, err := repo.GetAll(ctx)

	assert.NoError(t, err)
	assert.Equal(t, 1000, len(result))
}
func TestGenericRepository_Create_SpecialCharacters(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	model := mocks.TestModel{Email: "!@#$%^&*()_+-=~`"}
	_, err := repo.Create(ctx, model)
	assert.NoError(t, err)
}

func TestGenericRepository_Delete_Unscoped(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	model := mocks.TestModel{Email: "ToBeDeleted"}
	createdModel, err := repo.Create(ctx, model)
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}

	err = repo.Delete(ctx, createdModel.ID, true)
	if err != nil {
		t.Fatalf("Failed to delete model permanently: %v", err)
	}

	var result mocks.TestModel
	err = db.Unscoped().First(&result, createdModel.ID).Error
	if err == nil {
		t.Errorf("Expected record to be deleted, but it still exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Expected 'record not found' error, but got: %v", err)
	}
}

func TestGenericRepository_Delete_DBError(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	model := mocks.TestModel{Email: "error@test.com"}
	_, err := repo.Create(ctx, model)
	if err != nil {
		t.Fatalf("Failed to set up test data: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.Close()

	err = repo.Delete(ctx, model.ID, true)
	if err == nil {
		t.Fatalf("Expected an error due to DB failure, but got none")
	}
}
func TestGenericRepository_Delete_NonExistentItem(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	nonExistentID := uint(9999)

	err := repo.Delete(ctx, nonExistentID, true)
	if err == nil {
		t.Fatalf("Expected an error when deleting a non-existent item, but got none")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("Expected 'record not found' error, but got: %v", err)
	}
}

func TestGenericRepository_Get_DBError(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &mocks.TestModel{})
	repo := NewGenericRepository[mocks.TestModel, uint](db)

	model := mocks.TestModel{Email: "error@test.com"}
	_, err := repo.Create(ctx, model)
	if err != nil {
		t.Fatalf("Failed to create test model: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.Close()

	_, err = repo.Get(ctx, model.ID, "")
	if err == nil {
		t.Fatalf("Expected an error due to DB failure, but got none")
	}
}

func TestGenericRepository_Get_WithPreload(t *testing.T) {
	db := mocks.SetupGORMSqlite(t, &TestModelPreload{}, &OrdersModel{})
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

	result, err := repo.Get(ctx, user.ID, "Orders")
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
	db := mocks.SetupGORMSqlite(t, &TestModelWithStringID{})
	repo := NewGenericRepository[TestModelWithStringID, string](db)

	model := TestModelWithStringID{ID: "abc123", Email: "test@test.com"}
	err := db.Create(&model).Error
	if err != nil {
		t.Fatalf("Unexpected error creating model: %v", err)
	}

	err = repo.Delete(ctx, model.ID, false)
	if err != nil {
		t.Fatalf("Unexpected error during delete: %v", err)
	}

	var deletedModel TestModelWithStringID
	result := db.First(&deletedModel, "id = ?", model.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		t.Fatalf("Expected record to be deleted, but got unexpected error: %v", result.Error)
	}
}
