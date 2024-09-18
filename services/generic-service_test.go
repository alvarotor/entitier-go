package services

import (
	"errors"
	"testing"

	"github.com/alvarotor/entitier-go/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TestModel struct {
	ID    uint
	Name  string
	Email string
}

type MockRepo struct {
	mock.Mock
	repositories.IGenericRepo[TestModel, uint]
}

func (m *MockRepo) GetAll() ([]*TestModel, error) {
	args := m.Called()
	return args.Get(0).([]*TestModel), args.Error(1)
}

func (m *MockRepo) Get(ID uint, preload string) (*TestModel, error) {
	args := m.Called(ID, preload)
	return args.Get(0).(*TestModel), args.Error(1)
}

func (m *MockRepo) Create(data TestModel) (TestModel, error) {
	args := m.Called(data)
	return args.Get(0).(TestModel), args.Error(1)
}

func (m *MockRepo) Delete(ID uint, permanently bool) error {
	args := m.Called(ID, permanently)
	return args.Error(0)
}

func (m *MockRepo) Update(ID uint, amended TestModel) error {
	args := m.Called(ID, amended)
	return args.Error(0)
}

func TestNewGenericService(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewGenericService[TestModel, uint](mockRepo)
	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
}

func TestGetAll(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewGenericService[TestModel, uint](mockRepo)

	t.Run("Success", func(t *testing.T) {
		expected := []*TestModel{{ID: 1, Email: "test@test.com", Name: "Test"}}
		mockRepo.On("GetAll").Return(expected, nil).Once()

		result, err := service.GetAll()
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Empty Result", func(t *testing.T) {
		mockRepo.On("GetAll").Return([]*TestModel{}, nil).Once()

		result, err := service.GetAll()
		assert.NoError(t, err)
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		expectedError := errors.New("database error")
		mockRepo.On("GetAll").Return([]*TestModel(nil), expectedError).Once()

		result, err := service.GetAll()
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestGet(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewGenericService[TestModel, uint](mockRepo)

	t.Run("Success", func(t *testing.T) {
		expected := &TestModel{ID: 1, Name: "Test"}
		mockRepo.On("Get", uint(1), "").Return(expected, nil).Once()

		result, err := service.Get(1, "")
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		expectedError := errors.New("not found")
		mockRepo.On("Get", uint(2), "").Return((*TestModel)(nil), expectedError).Once()

		result, err := service.Get(2, "")
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreate(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewGenericService[TestModel, uint](mockRepo)

	t.Run("Success", func(t *testing.T) {
		input := TestModel{Name: "New Test"}
		expected := TestModel{ID: 1, Name: "New Test"}
		mockRepo.On("Create", input).Return(expected, nil).Once()

		result, err := service.Create(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		input := TestModel{Name: "Error Test"}
		expectedError := errors.New("creation error")
		mockRepo.On("Create", input).Return(TestModel{}, expectedError).Once()

		result, err := service.Create(input)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, input, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewGenericService[TestModel, uint](mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("Delete", uint(1), false).Return(nil).Once()

		err := service.Delete(1, false)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		expectedError := errors.New("deletion error")
		mockRepo.On("Delete", uint(2), true).Return(expectedError).Once()

		err := service.Delete(2, true)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockRepo := new(MockRepo)
	service := NewGenericService[TestModel, uint](mockRepo)

	t.Run("Success", func(t *testing.T) {
		amended := TestModel{ID: 1, Name: "Updated Test"}
		mockRepo.On("Update", uint(1), amended).Return(nil).Once()

		err := service.Update(1, amended)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		amended := TestModel{ID: 2, Name: "Error Test"}
		expectedError := errors.New("update error")
		mockRepo.On("Update", uint(2), amended).Return(expectedError).Once()

		err := service.Update(2, amended)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertExpectations(t)
	})
}
