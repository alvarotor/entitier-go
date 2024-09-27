package services_test

import (
	"errors"
	"testing"

	"github.com/alvarotor/entitier-go/mocks"
	"github.com/alvarotor/entitier-go/services"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGenericService_GetAll_Success(t *testing.T) {
	mockRepo := new(mocks.IGenericRepo[mocks.TestModel, uint])

	mockRepo.On("GetAll").Return([]*mocks.TestModel{
		{Email: "test1@example.com"},
		{Email: "test2@example.com"},
	}, nil)

	service := services.NewGenericService[mocks.TestModel, uint](mockRepo)
	result, err := service.GetAll()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "test1@example.com", result[0].Email)
}

func TestGenericService_GetAll_Empty(t *testing.T) {
	mockRepo := new(mocks.IGenericRepo[mocks.TestModel, uint])

	mockRepo.On("GetAll").Return([]*mocks.TestModel{}, nil)

	service := services.NewGenericService[mocks.TestModel, uint](mockRepo)
	result, err := service.GetAll()

	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestGenericService_GetAll_Error(t *testing.T) {
	mockRepo := new(mocks.IGenericRepo[mocks.TestModel, uint])

	mockRepo.On("GetAll").Return([]*mocks.TestModel{}, errors.New("DB error"))

	service := services.NewGenericService[mocks.TestModel, uint](mockRepo)
	result, err := service.GetAll()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "DB error", err.Error())
}

func TestGenericService_Get_Success(t *testing.T) {
	mockRepo := new(mocks.IGenericRepo[mocks.TestModel, uint])

	mockRepo.On("Get", uint(1), "").Return(&mocks.TestModel{ID: 1, Email: "test1@example.com"}, nil)

	service := services.NewGenericService[mocks.TestModel, uint](mockRepo)
	result, err := service.Get(1, "")

	assert.NoError(t, err)
	assert.Equal(t, "test1@example.com", result.Email)
}

func TestGenericService_Get_Error(t *testing.T) {
	mockRepo := new(mocks.IGenericRepo[mocks.TestModel, uint])

	mockRepo.On("Get", uint(1), "").Return(nil, gorm.ErrRecordNotFound)

	service := services.NewGenericService[mocks.TestModel, uint](mockRepo)
	result, err := service.Get(1, "")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestGenericService_Create_Success(t *testing.T) {
	mockRepo := new(mocks.IGenericRepo[mocks.TestModel, uint])

	inputModel := mocks.TestModel{Email: "test@example.com"}
	createdModel := mocks.TestModel{ID: 1, Email: "test@example.com"}

	mockRepo.On("Create", inputModel).Return(createdModel, nil)

	service := services.NewGenericService[mocks.TestModel, uint](mockRepo)

	result, err := service.Create(inputModel)

	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, createdModel, result)
}

func TestGenericService_Create_Error(t *testing.T) {
	mockRepo := new(mocks.IGenericRepo[mocks.TestModel, uint])

	inputModel := mocks.TestModel{Email: "test@example.com"}
	mockRepo.On("Create", inputModel).Return(mocks.TestModel{}, errors.New("DB error"))

	service := services.NewGenericService[mocks.TestModel, uint](mockRepo)

	result, err := service.Create(inputModel)

	assert.Error(t, err)
	assert.Equal(t, inputModel, result)
	assert.Equal(t, "DB error", err.Error())
}

func TestGenericService_Delete_Success(t *testing.T) {
	mockRepo := new(mocks.IGenericRepo[mocks.TestModel, uint])

	mockRepo.On("Delete", uint(1), true).Return(nil)

	service := services.NewGenericService[mocks.TestModel, uint](mockRepo)
	err := service.Delete(1, true)

	assert.NoError(t, err)
}

func TestGenericService_Delete_Error(t *testing.T) {
	mockRepo := new(mocks.IGenericRepo[mocks.TestModel, uint])

	mockRepo.On("Delete", uint(1), true).Return(gorm.ErrRecordNotFound)

	service := services.NewGenericService[mocks.TestModel, uint](mockRepo)
	err := service.Delete(1, true)

	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestGenericService_Update_Success(t *testing.T) {
	mockRepo := new(mocks.IGenericRepo[mocks.TestModel, uint])

	updatedModel := mocks.TestModel{Email: "updated@example.com"}

	mockRepo.On("Update", uint(1), updatedModel).Return(nil)

	service := services.NewGenericService[mocks.TestModel, uint](mockRepo)

	err := service.Update(1, updatedModel)

	assert.NoError(t, err)
	assert.Equal(t, "updated@example.com", updatedModel.Email)
}

func TestGenericService_Update_Error(t *testing.T) {
	mockRepo := new(mocks.IGenericRepo[mocks.TestModel, uint])

	updatedModel := mocks.TestModel{ID: 1, Email: "updated@example.com"}

	mockRepo.On("Update", uint(1), updatedModel).Return(gorm.ErrRecordNotFound)

	service := services.NewGenericService[mocks.TestModel, uint](mockRepo)

	err := service.Update(1, updatedModel)

	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
