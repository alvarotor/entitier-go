package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alvarotor/entitier-go/mocks"
	"github.com/alvarotor/entitier-go/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type TestModel struct {
	ID    uint   `gorm:"primaryKey"`
	Email string `gorm:"unique"`
}

// Mock Logger
type MockLogger struct{}

func (m *MockLogger) Info(msg string)  {}
func (m *MockLogger) Error(msg string) {}
func (m *MockLogger) Warn(msg string)  {}
func (m *MockLogger) Debug(msg string) {}
func (m *MockLogger) Fatal(msg string) {}
func (m *MockLogger) Trace(msg string) {}
func (m *MockLogger) Print(msg string) {}

// Mock Service
type MockService struct {
	mock.Mock
}

func (m *MockService) GetAll() ([]*TestModel, error) {
	args := m.Called()
	return args.Get(0).([]*TestModel), args.Error(1)
}

func (m *MockService) Get(ID uint, preload string) (*TestModel, error) {
	args := m.Called(ID, preload)
	return args.Get(0).(*TestModel), args.Error(1)
}

func (m *MockService) Create(data TestModel) (TestModel, error) {
	args := m.Called(data)
	return args.Get(0).(TestModel), args.Error(1)
}

func (m *MockService) Delete(ID uint, permanently bool) error {
	args := m.Called(ID, permanently)
	return args.Error(0)
}

func (m *MockService) Update(ID uint, amended TestModel) error {
	args := m.Called(ID, amended)
	return args.Error(0)
}

// Utility for creating mock gin context
func CreateMockGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestController_GetAll_Success(t *testing.T) {
	// Create a mock service
	mockService := new(mocks.IGenericService[TestModel, uint])
	mockLogger := &MockLogger{}

	// Prepare test data
	testModels := []*TestModel{
		{ID: 1, Email: "test1@example.com"},
		{ID: 2, Email: "test2@example.com"},
	}
	mockService.On("GetAll").Return(testModels, nil)

	// Create the controller
	ctrl := &controllerGeneric[TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	// Create a new Gin context for testing
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Call the controller method
	ctrl.GetAll(c)

	// Assert that the response status code is 200
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert that the response contains the expected data
	expectedBody := `{"all":[{"ID":1,"Email":"test1@example.com"},{"ID":2,"Email":"test2@example.com"}]}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestController_GetAll_NotFound(t *testing.T) {
	// Create a mock service
	mockService := new(mocks.IGenericService[TestModel, uint])
	mockLogger := &MockLogger{}

	// Simulate a "not found" error
	mockService.On("GetAll").Return(nil, models.ErrNotFound)

	// Create the controller
	ctrl := &controllerGeneric[TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	// Create a new Gin context for testing
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Call the controller method
	ctrl.GetAll(c)

	// Assert that the response status code is 404 (Not Found)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Assert that the response contains the expected error message
	expectedBody := fmt.Sprintf(`{"err":"%s"}`, models.ErrNotFound.Error())
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestController_GetAll_InternalError(t *testing.T) {
	// Create a mock service
	mockService := new(mocks.IGenericService[TestModel, uint])
	mockLogger := &MockLogger{}

	// Simulate a generic internal error
	mockService.On("GetAll").Return(nil, errors.New("database error"))

	// Create the controller
	ctrl := &controllerGeneric[TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	// Create a new Gin context for testing
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Call the controller method
	ctrl.GetAll(c)

	// Assert that the response status code is 500 (Internal Server Error)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Assert that the response contains the expected error message
	expectedBody := `{"err":"database error"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

// --- Test Create ---

func TestController_Create_Success(t *testing.T) {
	// Create a mock service
	mockService := new(mocks.IGenericService[TestModel, uint])
	mockLogger := &MockLogger{}

	// Test data
	inputModel := TestModel{Email: "test@example.com"}
	createdModel := TestModel{ID: 1, Email: "test@example.com"}

	// Simulate successful creation
	mockService.On("Create", inputModel).Return(createdModel, nil)

	// Create the controller
	ctrl := &controllerGeneric[TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	// Call the controller method
	result, err := ctrl.Create(inputModel)

	// Assert the result
	assert.NoError(t, err)
	assert.Equal(t, createdModel, result)
}

func TestController_Create_Failure(t *testing.T) {
	// Create a mock service
	mockService := new(mocks.IGenericService[TestModel, uint])
	mockLogger := &MockLogger{}

	// Test data
	inputModel := TestModel{Email: "test@example.com"}

	// Simulate creation failure
	mockService.On("Create", inputModel).Return(inputModel, errors.New("creation error"))

	// Create the controller
	ctrl := &controllerGeneric[TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	// Call the controller method
	result, err := ctrl.Create(inputModel)

	// Assert the result
	assert.Error(t, err)
	assert.Equal(t, inputModel, result)
	assert.Equal(t, "creation error", err.Error())
}

func TestController_Get_Success(t *testing.T) {
	// Create a mock service
	mockService := new(mocks.IGenericService[TestModel, uint])
	mockLogger := &MockLogger{}

	// Prepare test data
	testModel := &TestModel{ID: 1, Email: "test1@example.com"}
	mockService.On("Get", uint(1), "User").Return(testModel, nil)

	// Create the controller
	ctrl := &controllerGeneric[TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	// Create a new Gin context for testing
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Simulate the presence of the ID in the URL parameters to simulate utils.GetIDParam
	c.Params = gin.Params{
		{Key: "id", Value: "1"}, // Simulate ID param from utils.GetIDParam
	}

	// Call the controller method
	ctrl.Get(c)

	// Assert that the response status code is 200
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert that the response contains the expected data
	expectedBody := `{"item":{"ID":1,"Email":"test1@example.com"}}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestController_Get_NotFoundVariants(t *testing.T) {
	tests := []struct {
		name         string
		mockError    error
		expectedBody string
		expectedCode int
	}{
		{"Model Not Found", models.ErrNotFound, `{"err":"` + models.ErrNotFound.Error() + `"}`, http.StatusInternalServerError},
		{"GORM Record Not Found", gorm.ErrRecordNotFound, `{"err":"no rows found"}`, http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock service
			mockService := new(mocks.IGenericService[TestModel, uint])
			mockLogger := &mocks.Logger{}

			// Simulate the not found errors
			mockService.On("Get", uint(1), "User").Return(nil, tt.mockError)

			// Create the controller
			ctrl := &controllerGeneric[TestModel, uint]{
				svcT: mockService,
				log:  mockLogger,
			}

			// Create a new Gin context for testing
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: "1"}}

			// Call the controller method
			ctrl.Get(c)

			// Assert the response code and body
			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestController_Delete_Success(t *testing.T) {
	// Create a mock service
	mockService := new(mocks.IGenericService[TestModel, uint])
	mockLogger := &MockLogger{}

	// Simulate successful deletion
	mockService.On("Delete", uint(1), true).Return(nil)

	// Create the controller
	ctrl := &controllerGeneric[TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	// Create a new Gin context for testing
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Call the controller method
	ctrl.Delete(c)

	// Assert that the response status code is 200 (OK)
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert that the response contains the success message
	expectedBody := `{"message":"deleted"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestController_Delete_Failure(t *testing.T) {
	// Create a mock service
	mockService := new(mocks.IGenericService[TestModel, uint])
	mockLogger := &MockLogger{}

	// Simulate deletion failure
	mockService.On("Delete", uint(1), true).Return(errors.New("deletion error"))

	// Create the controller
	ctrl := &controllerGeneric[TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	// Create a new Gin context for testing
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Call the controller method
	ctrl.Delete(c)

	// Assert that the response status code is 500 (Internal Server Error)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Assert that the response contains the expected error message
	expectedBody := `{"err":"deletion error"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestController_Update_Success(t *testing.T) {
	// Create a mock service
	mockService := new(mocks.IGenericService[TestModel, uint])
	mockLogger := &MockLogger{}

	// Test data
	inputModel := TestModel{ID: 1, Email: "updated@example.com"}

	// Simulate successful update
	mockService.On("Update", uint(1), inputModel).Return(nil)

	// Create the controller
	ctrl := &controllerGeneric[TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	// Create a new Gin context for testing
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Call the controller method
	ctrl.Update(c, inputModel)

	// Assert that the response status code is 200 (OK)
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert that the response contains the success message
	expectedBody := `{"message": "updated"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestController_Update_Failure(t *testing.T) {
	// Create a mock service
	mockService := new(mocks.IGenericService[TestModel, uint])
	mockLogger := &MockLogger{}

	// Test data
	inputModel := TestModel{ID: 1, Email: "updated@example.com"}

	// Simulate update failure
	mockService.On("Update", uint(1), inputModel).Return(errors.New("update error"))

	// Create the controller
	ctrl := &controllerGeneric[TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	// Create a new Gin context for testing
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Call the controller method
	ctrl.Update(c, inputModel)

	// Assert that the response status code is 500 (Internal Server Error)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Assert that the response contains the expected error message
	expectedBody := `{"err":"update error"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestController_Get_IDParamInvalid(t *testing.T) {
	// Create a mock service
	mockService := new(mocks.IGenericService[TestModel, uint])
	mockLogger := &MockLogger{}

	// Create the controller
	ctrl := &controllerGeneric[TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	// Create a new Gin context for testing
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Simulate an invalid ID in the URL parameters
	c.Params = gin.Params{
		{Key: "id", Value: "abc"}, // Simulate invalid ID
	}

	// Call the controller method
	ctrl.Get(c)

	// Assert that the response status code is 400 (Bad Request)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Assert that the response contains the expected error message
	expectedBody := `{"err":"id type mismatch"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestController_Get_InternalError(t *testing.T) {
	// Create a mock service
	mockService := new(mocks.IGenericService[TestModel, uint])
	mockLogger := &MockLogger{}

	// Simulate an internal error in the mock service
	mockService.On("Get", mock.Anything, "User").Return(nil, errors.New("some internal error"))

	// Create the controller
	ctrl := &controllerGeneric[TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	// Create a new Gin context for testing
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Simulate the presence of the ID in the URL parameters
	c.Params = gin.Params{
		{Key: "id", Value: "1"},
	}

	// Call the controller's Get method
	ctrl.Get(c)

	// Assert that the response status code is 500 (Internal Server Error)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Assert that the response contains the expected error message
	expectedBody := `{"err":"some internal error"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}
func TestController_Get_InvalidOrNilParam(t *testing.T) {
	tests := []struct {
		name         string
		paramKey     string
		paramValue   string
		expectedCode int
		expectedBody string
	}{
		{"Missing ID Param", "", "", http.StatusBadRequest, `{"err":"must provide valid id"}`},
		{"Invalid ID Param", "id", "abc", http.StatusBadRequest, `{"err":"id type mismatch"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock service
			mockService := new(mocks.IGenericService[TestModel, uint])
			mockLogger := &MockLogger{}

			// Create the controller
			ctrl := &controllerGeneric[TestModel, uint]{
				svcT: mockService,
				log:  mockLogger,
			}

			// Create a new Gin context for testing
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Conditionally set the parameter in the Gin context
			if tt.paramKey != "" {
				c.Params = gin.Params{{Key: tt.paramKey, Value: tt.paramValue}}
			}

			// Call the controller method
			ctrl.Get(c)

			// Assert the response code
			assert.Equal(t, tt.expectedCode, w.Code)

			// Assert that the response contains the expected error message
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
func TestController_Update_InvalidOrNilParam(t *testing.T) {
	tests := []struct {
		name         string
		paramKey     string
		paramValue   string
		expectedCode int
		expectedBody string
	}{
		{"Missing ID Param", "", "", http.StatusBadRequest, `{"err":"must provide valid id"}`},
		{"Invalid ID Param", "id", "abc", http.StatusBadRequest, `{"err":"id type mismatch"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock service
			mockService := new(mocks.IGenericService[TestModel, uint])
			mockLogger := &MockLogger{}

			// Create the controller
			ctrl := &controllerGeneric[TestModel, uint]{
				svcT: mockService,
				log:  mockLogger,
			}

			// Create a new Gin context for testing
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Conditionally set the parameter in the Gin context
			if tt.paramKey != "" {
				c.Params = gin.Params{{Key: tt.paramKey, Value: tt.paramValue}}
			}

			// Call the controller method
			ctrl.Update(c, TestModel{})

			// Assert the response code
			assert.Equal(t, tt.expectedCode, w.Code)

			// Assert that the response contains the expected error message
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
func TestController_DeleteInvalidOrNilParam(t *testing.T) {
	tests := []struct {
		name         string
		paramKey     string
		paramValue   string
		expectedCode int
		expectedBody string
	}{
		{"Missing ID Param", "", "", http.StatusBadRequest, `{"err":"must provide valid id"}`},
		{"Invalid ID Param", "id", "abc", http.StatusBadRequest, `{"err":"id type mismatch"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock service
			mockService := new(mocks.IGenericService[TestModel, uint])
			mockLogger := &MockLogger{}

			// Create the controller
			ctrl := &controllerGeneric[TestModel, uint]{
				svcT: mockService,
				log:  mockLogger,
			}

			// Create a new Gin context for testing
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Conditionally set the parameter in the Gin context
			if tt.paramKey != "" {
				c.Params = gin.Params{{Key: tt.paramKey, Value: tt.paramValue}}
			}

			// Call the controller method
			ctrl.Delete(c)

			// Assert the response code
			assert.Equal(t, tt.expectedCode, w.Code)

			// Assert that the response contains the expected error message
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
func TestController_Update_Failure_Saving(t *testing.T) {
	// Create a mock service
	mockService := new(mocks.IGenericService[TestModel, uint])
	mockLogger := &MockLogger{}

	// Test data
	inputModel := TestModel{ID: 1, Email: "updated@example.com"}

	// Simulate update failure
	mockService.On("Update", uint(1), inputModel).Return(gorm.ErrRecordNotFound)

	// Create the controller
	ctrl := &controllerGeneric[TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	// Create a new Gin context for testing
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	// Call the controller method
	ctrl.Update(c, inputModel)

	// Assert that the response status code is 500 (Internal Server Error)
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Assert that the response contains the expected error message
	expectedBody := `{"err":"no rows found"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

// assert.ErrorIs(t, err, gorm.ErrRecordNotFound, "Expected gorm.ErrRecordNotFound error")
