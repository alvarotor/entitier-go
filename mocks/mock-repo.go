package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
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
