package services_test

import (
	"errors"
	"testing"

	"github.com/alvarotor/entitier-go/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type TestModel struct {
	ID    uint   `gorm:"primaryKey"`
	Email string `gorm:"unique"`
}

// Mock Repo Struct
type MockRepo struct {
	data        []*TestModel
	errOnGet    error
	errOnGetAll error
	errOnCreate error
	errOnDelete error
	errOnUpdate error
}

func (m *MockRepo) GetAll() ([]*TestModel, error) {
	if m.errOnGetAll != nil {
		return nil, m.errOnGetAll
	}
	return m.data, nil
}

func (m *MockRepo) Get(ID uint, preload string) (*TestModel, error) {
	if m.errOnGet != nil {
		return nil, m.errOnGet
	}
	for _, model := range m.data {
		if model.ID == ID {
			return model, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockRepo) Create(data TestModel) (TestModel, error) {
	if m.errOnCreate != nil {
		return data, m.errOnCreate
	}
	m.data = append(m.data, &data)
	return data, nil
}

func (m *MockRepo) Delete(ID uint, permanently bool) error {
	if m.errOnDelete != nil {
		return m.errOnDelete
	}
	for i, model := range m.data {
		if model.ID == ID {
			m.data = append(m.data[:i], m.data[i+1:]...)
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *MockRepo) Update(ID uint, amended TestModel) error {
	if m.errOnUpdate != nil {
		return m.errOnUpdate
	}
	for _, model := range m.data {
		if model.ID == ID {
			model.Email = amended.Email
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

// Test file

func TestGenericService_GetAll_Success(t *testing.T) {
	mockRepo := &MockRepo{
		data: []*TestModel{
			{ID: 1, Email: "test1@example.com"},
			{ID: 2, Email: "test2@example.com"},
		},
	}

	service := services.NewGenericService[TestModel, uint](mockRepo)
	result, err := service.GetAll()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "test1@example.com", result[0].Email)
}

func TestGenericService_GetAll_Empty(t *testing.T) {
	mockRepo := &MockRepo{
		data: []*TestModel{},
	}

	service := services.NewGenericService[TestModel, uint](mockRepo)
	result, err := service.GetAll()

	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestGenericService_GetAll_Error(t *testing.T) {
	mockRepo := &MockRepo{
		errOnGetAll: errors.New("DB error"),
	}

	service := services.NewGenericService[TestModel, uint](mockRepo)
	result, err := service.GetAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "DB error", err.Error())
}

func TestGenericService_Get_Success(t *testing.T) {
	mockRepo := &MockRepo{
		data: []*TestModel{
			{ID: 1, Email: "test1@example.com"},
		},
	}

	service := services.NewGenericService[TestModel, uint](mockRepo)
	result, err := service.Get(1, "")

	assert.NoError(t, err)
	assert.Equal(t, "test1@example.com", result.Email)
}

func TestGenericService_Get_Error(t *testing.T) {
	mockRepo := &MockRepo{
		errOnGet: gorm.ErrRecordNotFound,
	}

	service := services.NewGenericService[TestModel, uint](mockRepo)
	result, err := service.Get(1, "")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestGenericService_Create_Success(t *testing.T) {
	mockRepo := &MockRepo{}
	service := services.NewGenericService[TestModel, uint](mockRepo)

	model := TestModel{Email: "test1@example.com"}
	result, err := service.Create(model)

	assert.NoError(t, err)
	assert.Equal(t, "test1@example.com", result.Email)
	assert.Equal(t, 1, len(mockRepo.data))
}

func TestGenericService_Create_Error(t *testing.T) {
	mockRepo := &MockRepo{
		errOnCreate: errors.New("DB error"),
	}

	service := services.NewGenericService[TestModel, uint](mockRepo)

	model := TestModel{Email: "test1@example.com"}
	result, err := service.Create(model)

	assert.Error(t, err)
	assert.Equal(t, model, result)
	assert.Equal(t, "DB error", err.Error())
}

func TestGenericService_Delete_Success(t *testing.T) {
	mockRepo := &MockRepo{
		data: []*TestModel{
			{ID: 1, Email: "test1@example.com"},
		},
	}

	service := services.NewGenericService[TestModel, uint](mockRepo)
	err := service.Delete(1, true)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(mockRepo.data))
}

func TestGenericService_Delete_Error(t *testing.T) {
	mockRepo := &MockRepo{
		errOnDelete: gorm.ErrRecordNotFound,
	}

	service := services.NewGenericService[TestModel, uint](mockRepo)
	err := service.Delete(1, true)

	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestGenericService_Update_Success(t *testing.T) {
	mockRepo := &MockRepo{
		data: []*TestModel{
			{ID: 1, Email: "test1@example.com"},
		},
	}

	service := services.NewGenericService[TestModel, uint](mockRepo)

	updatedModel := TestModel{Email: "updated@example.com"}
	err := service.Update(1, updatedModel)

	assert.NoError(t, err)
	assert.Equal(t, "updated@example.com", mockRepo.data[0].Email)
}

func TestGenericService_Update_Error(t *testing.T) {
	mockRepo := &MockRepo{
		errOnUpdate: gorm.ErrRecordNotFound,
	}

	service := services.NewGenericService[TestModel, uint](mockRepo)

	updatedModel := TestModel{Email: "updated@example.com"}
	err := service.Update(1, updatedModel)

	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
