package mocks

import "github.com/stretchr/testify/mock"

// MockLogger is a mock implementation of the logger.Logger interface
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string) {
	m.Called(msg)
}

func (m *MockLogger) Error(msg string) {
	m.Called(msg)
}

func (m *MockLogger) Debug(msg string) {
	m.Called(msg)
}
