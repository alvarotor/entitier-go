package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockRepo is a mock implementation of IGenericRepo
type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Create(data MockModel) (MockModel, error) {
	args := m.Called(data)
	return args.Get(0).(MockModel), args.Error(1)
}

func (m *MockRepo) GetAll() ([]*MockModel, error) {
	args := m.Called()
	return args.Get(0).([]*MockModel), args.Error(1)
}

func (m *MockRepo) Get(id uint, preload string) (*MockModel, error) {
	args := m.Called(id, preload)
	return args.Get(0).(*MockModel), args.Error(1)
}

func (m *MockRepo) Update(id uint, data MockModel) error {
	args := m.Called(id, data)
	return args.Error(0)
}

func (m *MockRepo) Delete(id uint, permanently bool) error {
	args := m.Called(id, permanently)
	return args.Error(0)
}
