package services

// import (
// 	"errors"
// 	"testing"

// 	"github.com/alvarotor/entitier-go/mocks"
// 	"github.com/stretchr/testify/assert"
// )

// func TestNewGenericService(t *testing.T) {
// 	mockRepo := new(mocks.MockRepo)
// 	service := NewGenericService[mocks.MockModel, uint](mockRepo)

// 	assert.NotNil(t, service)
// 	assert.IsType(t, &GenericService[mocks.MockModel, uint]{}, service)
// }

// func TestGenericService_GetAll(t *testing.T) {
// 	mockRepo := new(mocks.MockRepo)
// 	service := NewGenericService[mocks.MockModel, uint](mockRepo)

// 	t.Run("Success", func(t *testing.T) {
// 		expected := []*mocks.MockModel{{ID: 1, Email: "test@example.com", Name: "Test"}}
// 		mockRepo.On("GetAll").Return(expected, nil).Once()

// 		result, err := service.GetAll()

// 		assert.NoError(t, err)
// 		assert.Equal(t, expected, result)
// 		mockRepo.AssertExpectations(t)
// 	})

// 	t.Run("Empty Result", func(t *testing.T) {
// 		mockRepo.On("GetAll").Return([]*mocks.MockModel{}, nil).Once()

// 		result, err := service.GetAll()

// 		assert.NoError(t, err)
// 		assert.Empty(t, result)
// 		mockRepo.AssertExpectations(t)
// 	})

// 	t.Run("Error", func(t *testing.T) {
// 		expectedError := errors.New("database error")
// 		mockRepo.On("GetAll").Return([]*mocks.MockModel(nil), expectedError).Once()

// 		result, err := service.GetAll()

// 		assert.Error(t, err)
// 		assert.Equal(t, expectedError, err)
// 		assert.Nil(t, result)
// 		mockRepo.AssertExpectations(t)
// 	})
// }

// func TestGenericService_Get(t *testing.T) {
// 	mockRepo := new(mocks.MockRepo)
// 	service := NewGenericService[mocks.MockModel, uint](mockRepo)

// 	t.Run("Success", func(t *testing.T) {
// 		expected := &mocks.MockModel{ID: 1, Email: "test@example.com", Name: "Test"}
// 		mockRepo.On("Get", uint(1), "").Return(expected, nil).Once()

// 		result, err := service.Get(1, "")

// 		assert.NoError(t, err)
// 		assert.Equal(t, expected, result)
// 		mockRepo.AssertExpectations(t)
// 	})

// 	t.Run("Error", func(t *testing.T) {
// 		expectedError := errors.New("not found")
// 		mockRepo.On("Get", uint(2), "").Return((*mocks.MockModel)(nil), expectedError).Once()

// 		result, err := service.Get(2, "")

// 		assert.Error(t, err)
// 		assert.Equal(t, expectedError, err)
// 		assert.Nil(t, result)
// 		mockRepo.AssertExpectations(t)
// 	})
// }

// func TestGenericService_Create(t *testing.T) {
// 	mockRepo := new(mocks.MockRepo)
// 	service := NewGenericService[mocks.MockModel, uint](mockRepo)

// 	t.Run("Success", func(t *testing.T) {
// 		input := mocks.MockModel{Email: "new@example.com", Name: "New"}
// 		expected := mocks.MockModel{ID: 1, Email: "new@example.com", Name: "New"}
// 		mockRepo.On("Create", input).Return(expected, nil).Once()

// 		result, err := service.Create(input)

// 		assert.NoError(t, err)
// 		assert.Equal(t, expected, result)
// 		mockRepo.AssertExpectations(t)
// 	})

// 	t.Run("Error", func(t *testing.T) {
// 		input := mocks.MockModel{Email: "invalid@example.com", Name: "Invalid"}
// 		expectedError := errors.New("creation failed")
// 		mockRepo.On("Create", input).Return(mocks.MockModel{}, expectedError).Once()

// 		result, err := service.Create(input)

// 		assert.Error(t, err)
// 		assert.Equal(t, expectedError, err)
// 		assert.Equal(t, input, result)
// 		mockRepo.AssertExpectations(t)
// 	})
// }

// func TestGenericService_Delete(t *testing.T) {
// 	mockRepo := new(mocks.MockRepo)
// 	service := NewGenericService[mocks.MockModel, uint](mockRepo)

// 	t.Run("Success", func(t *testing.T) {
// 		mockRepo.On("Delete", uint(1), false).Return(nil).Once()

// 		err := service.Delete(1, false)

// 		assert.NoError(t, err)
// 		mockRepo.AssertExpectations(t)
// 	})

// 	t.Run("Error", func(t *testing.T) {
// 		expectedError := errors.New("deletion failed")
// 		mockRepo.On("Delete", uint(2), true).Return(expectedError).Once()

// 		err := service.Delete(2, true)

// 		assert.Error(t, err)
// 		assert.Equal(t, expectedError, err)
// 		mockRepo.AssertExpectations(t)
// 	})
// }

// func TestGenericService_Update(t *testing.T) {
// 	mockRepo := new(mocks.MockRepo)
// 	service := NewGenericService[mocks.MockModel, uint](mockRepo)

// 	t.Run("Success", func(t *testing.T) {
// 		updatedModel := mocks.MockModel{ID: 1, Email: "updated@example.com", Name: "Updated"}
// 		mockRepo.On("Update", uint(1), updatedModel).Return(nil).Once()

// 		err := service.Update(1, updatedModel)

// 		assert.NoError(t, err)
// 		mockRepo.AssertExpectations(t)
// 	})

// 	t.Run("Error", func(t *testing.T) {
// 		updatedModel := mocks.MockModel{ID: 2, Email: "fail@example.com", Name: "Fail"}
// 		expectedError := errors.New("update failed")
// 		mockRepo.On("Update", uint(2), updatedModel).Return(expectedError).Once()

// 		err := service.Update(2, updatedModel)

// 		assert.Error(t, err)
// 		assert.Equal(t, expectedError, err)
// 		mockRepo.AssertExpectations(t)
// 	})
// }
