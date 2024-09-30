package mocks

import (
	gin "github.com/gin-gonic/gin"
	mock "github.com/stretchr/testify/mock"
)

type MockUtils[X string | uint] struct {
	mock.Mock
}

func (m *MockUtils[X]) GetIDParam(c *gin.Context) interface{} {
	args := m.Called(c)
	return args.Get(0)
}

func (m *MockUtils[X]) ConvertToGenericID(idInterface interface{}) (X, error) {
	args := m.Called(idInterface)
	return args.Get(0).(X), args.Error(1)
}
