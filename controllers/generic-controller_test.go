package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
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

var ctx = context.Background()

func CreateMockGinContext() (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestController_GetAll_Success(t *testing.T) {
	mockService := new(mocks.IGenericService[mocks.TestModel, uint])
	mockLogger := &mocks.Logger{}

	testModels := []*mocks.TestModel{
		{ID: 1, Email: "test1@example.com"},
		{ID: 2, Email: "test2@example.com"},
	}

	ctrl := &controllerGeneric[mocks.TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	c, w := CreateMockGinContext()

	mockService.On("GetAll", c).Return(testModels, nil)

	ctrl.GetAll(c)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedBody := `{"all":[{"ID":1,"Email":"test1@example.com"},{"ID":2,"Email":"test2@example.com"}]}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestController_GetAll_NotFound(t *testing.T) {
	mockService := new(mocks.IGenericService[mocks.TestModel, uint])
	mockLogger := &mocks.Logger{}

	mockLogger.On("Error", "no rows found").Return(nil)

	ctrl := &controllerGeneric[mocks.TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	c, w := CreateMockGinContext()

	mockService.On("GetAll", c).Return(nil, models.ErrNotFound)

	ctrl.GetAll(c)

	assert.Equal(t, http.StatusNotFound, w.Code)

	expectedBody := fmt.Sprintf(`{"err":"%s"}`, models.ErrNotFound.Error())
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestController_GetAll_InternalError(t *testing.T) {
	mockService := new(mocks.IGenericService[mocks.TestModel, uint])
	mockLogger := &mocks.Logger{}

	err := errors.New("database error")
	mockLogger.On("Error", err.Error()).Return(nil)

	ctrl := &controllerGeneric[mocks.TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	c, w := CreateMockGinContext()

	mockService.On("GetAll", c).Return(nil, err)

	ctrl.GetAll(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expectedBody := fmt.Sprintf(`{"err":"%s"}`, err.Error())
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestController_Create_Success(t *testing.T) {
	mockService := new(mocks.IGenericService[mocks.TestModel, uint])
	mockLogger := &mocks.Logger{}

	inputModel := mocks.TestModel{Email: "test@example.com"}
	createdModel := mocks.TestModel{ID: 1, Email: "test@example.com"}

	mockService.On("Create", ctx, inputModel).Return(createdModel, nil)

	ctrl := &controllerGeneric[mocks.TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	result, err := ctrl.Create(ctx, inputModel)

	assert.NoError(t, err)
	assert.Equal(t, createdModel, result)
}

func TestController_Create_Failure(t *testing.T) {
	// Create a mock service
	mockService := new(mocks.IGenericService[mocks.TestModel, uint])
	mockLogger := &mocks.Logger{}

	// Test data
	inputModel := mocks.TestModel{Email: "test@example.com"}

	// Simulate creation failure
	err := errors.New("create error")
	mockLogger.On("Error", err.Error()).Return(nil)
	mockService.On("Create", ctx, inputModel).Return(inputModel, err)

	// Create the controller
	ctrl := &controllerGeneric[mocks.TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	// Call the controller method
	result, err := ctrl.Create(ctx, inputModel)

	// Assert the result
	assert.Error(t, err)
	assert.Equal(t, inputModel, result)
	assert.Equal(t, err.Error(), err.Error())
}

func TestController_Get_Success(t *testing.T) {
	mockService := new(mocks.IGenericService[mocks.TestModel, uint])
	mockLogger := &mocks.Logger{}

	testModel := &mocks.TestModel{ID: 1, Email: "test1@example.com"}

	ctrl := &controllerGeneric[mocks.TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	c, w := CreateMockGinContext()

	c.Params = gin.Params{
		{Key: "id", Value: "1"},
	}

	mockService.On("Get", c, uint(1), "User").Return(testModel, nil)

	ctrl.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)

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
			mockService := new(mocks.IGenericService[mocks.TestModel, uint])
			mockLogger := &mocks.Logger{}

			ctrl := &controllerGeneric[mocks.TestModel, uint]{
				svcT: mockService,
				log:  mockLogger,
			}

			c, w := CreateMockGinContext()

			c.Params = gin.Params{{Key: "id", Value: "1"}}

			mockService.On("Get", c, uint(1), "User").Return(nil, tt.mockError)

			ctrl.Get(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestController_Delete_Success(t *testing.T) {
	mockService := new(mocks.IGenericService[mocks.TestModel, uint])
	mockLogger := &mocks.Logger{}

	ctrl := &controllerGeneric[mocks.TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	c, w := CreateMockGinContext()

	c.Params = gin.Params{{Key: "id", Value: "1"}}

	mockService.On("Delete", c, uint(1), true).Return(nil)

	ctrl.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedBody := `{"message":"deleted"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestController_Delete_Failure(t *testing.T) {
	mockService := new(mocks.IGenericService[mocks.TestModel, uint])
	mockLogger := &mocks.Logger{}

	mockLogger.On("Error", models.ErrNotFound.Error()).Return(nil)

	ctrl := &controllerGeneric[mocks.TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	c, w := CreateMockGinContext()

	c.Params = gin.Params{{Key: "id", Value: "1"}}

	mockService.On("Delete", c, uint(1), true).Return(models.ErrNotFound)

	ctrl.Delete(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expectedBody := fmt.Sprintf(`{"err":"%s"}`, models.ErrNotFound.Error())
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestController_Update_Success(t *testing.T) {
	mockService := new(mocks.IGenericService[mocks.TestModel, uint])
	mockLogger := &mocks.Logger{}

	inputModel := mocks.TestModel{ID: 1, Email: "updated@example.com"}

	ctrl := &controllerGeneric[mocks.TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	c, w := CreateMockGinContext()

	c.Params = gin.Params{{Key: "id", Value: "1"}}

	mockService.On("Update", c, uint(1), inputModel).Return(nil)

	ctrl.Update(c, inputModel)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedBody := `{"message": "updated"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestController_Update_Failure(t *testing.T) {
	mockService := new(mocks.IGenericService[mocks.TestModel, uint])
	mockLogger := &mocks.Logger{}

	inputModel := mocks.TestModel{ID: 1, Email: "updated@example.com"}

	ctrl := &controllerGeneric[mocks.TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	c, w := CreateMockGinContext()

	c.Params = gin.Params{{Key: "id", Value: "1"}}

	mockService.On("Update", c, uint(1), inputModel).Return(errors.New("update error"))

	ctrl.Update(c, inputModel)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expectedBody := `{"err":"update error"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestController_Get_IDParamInvalid(t *testing.T) {
	mockService := new(mocks.IGenericService[mocks.TestModel, uint])
	mockLogger := &mocks.Logger{}

	mockLogger.On("Error", models.ErrIDTypeMismatch.Error()).Return(nil)

	ctrl := &controllerGeneric[mocks.TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	c, w := CreateMockGinContext()

	c.Params = gin.Params{
		{Key: "id", Value: "abc"},
	}

	ctrl.Get(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expectedBody := fmt.Sprintf(`{"err":"%s"}`, models.ErrIDTypeMismatch.Error())
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestController_Get_InternalError(t *testing.T) {
	mockService := new(mocks.IGenericService[mocks.TestModel, uint])
	mockLogger := &mocks.Logger{}

	ctrl := &controllerGeneric[mocks.TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	c, w := CreateMockGinContext()

	c.Params = gin.Params{
		{Key: "id", Value: "1"},
	}

	mockService.On("Get", c, mock.Anything, "User").Return(nil, errors.New("some internal error"))

	ctrl.Get(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expectedBody := `{"err":"some internal error"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestController_Get_InvalidOrNilParam(t *testing.T) {
	tests := []struct {
		name             string
		paramKey         string
		paramValue       string
		expectedCode     int
		expectedBody     string
		expectedLogError string
	}{
		{"Missing ID Param", "", "", http.StatusBadRequest, fmt.Sprintf(`{"err":"%s"}`, models.ErrMustProvideValidID.Error()), models.ErrMustProvideValidID.Error()},
		{"Invalid ID Param", "id", "abc", http.StatusBadRequest, fmt.Sprintf(`{"err":"%s"}`, models.ErrIDTypeMismatch.Error()), models.ErrIDTypeMismatch.Error()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock service
			mockService := new(mocks.IGenericService[mocks.TestModel, uint])
			mockLogger := &mocks.Logger{}

			mockLogger.On("Error", tt.expectedLogError).Return(nil)

			// Create the controller
			ctrl := &controllerGeneric[mocks.TestModel, uint]{
				svcT: mockService,
				log:  mockLogger,
			}

			c, w := CreateMockGinContext()

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
		name             string
		paramKey         string
		paramValue       string
		expectedCode     int
		expectedBody     string
		expectedLogError string
	}{
		{"Missing ID Param", "", "", http.StatusBadRequest, fmt.Sprintf(`{"err":"%s"}`, models.ErrMustProvideValidID), models.ErrMustProvideValidID.Error()},
		{"Invalid ID Param", "id", "abc", http.StatusBadRequest, fmt.Sprintf(`{"err":"%s"}`, models.ErrIDTypeMismatch), models.ErrIDTypeMismatch.Error()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock service
			mockService := new(mocks.IGenericService[mocks.TestModel, uint])
			mockLogger := &mocks.Logger{}

			mockLogger.On("Error", tt.expectedLogError).Return(nil)

			// Create the controller
			ctrl := &controllerGeneric[mocks.TestModel, uint]{
				svcT: mockService,
				log:  mockLogger,
			}

			c, w := CreateMockGinContext()

			// Conditionally set the parameter in the Gin context
			if tt.paramKey != "" {
				c.Params = gin.Params{{Key: tt.paramKey, Value: tt.paramValue}}
			}

			// Call the controller method
			ctrl.Update(c, mocks.TestModel{})

			// Assert the response code
			assert.Equal(t, tt.expectedCode, w.Code)

			// Assert that the response contains the expected error message
			log.Println(tt.expectedBody)
			log.Println(w.Body.String())
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestController_Delete_InvalidOrNilParam(t *testing.T) {
	tests := []struct {
		name             string
		paramKey         string
		paramValue       string
		expectedCode     int
		expectedBody     string
		expectedLogError string
	}{
		{"Missing ID Param", "", "", http.StatusBadRequest, fmt.Sprintf(`{"err":"%s"}`, models.ErrMustProvideValidID.Error()), models.ErrMustProvideValidID.Error()},
		{"Invalid ID Param", "id", "abc", http.StatusBadRequest, fmt.Sprintf(`{"err":"%s"}`, models.ErrIDTypeMismatch.Error()), models.ErrIDTypeMismatch.Error()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock service
			mockService := new(mocks.IGenericService[mocks.TestModel, uint])
			mockLogger := &mocks.Logger{}

			mockLogger.On("Error", tt.expectedLogError).Return(nil)

			// Create the controller
			ctrl := &controllerGeneric[mocks.TestModel, uint]{
				svcT: mockService,
				log:  mockLogger,
			}

			c, w := CreateMockGinContext()

			// Conditionally set the parameter in the Gin context
			if tt.paramKey != "" {
				c.Params = gin.Params{{Key: tt.paramKey, Value: tt.paramValue}}
			}

			// Call the controller method
			ctrl.Delete(c)

			// Assert the response code
			assert.Equal(t, tt.expectedCode, w.Code)

			// Assert that the response contains the expected error message
			log.Println(tt.expectedBody)
			log.Println(w.Body.String())
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestController_Update_Failure_Saving(t *testing.T) {
	mockService := new(mocks.IGenericService[mocks.TestModel, uint])
	mockLogger := &mocks.Logger{}

	inputModel := mocks.TestModel{ID: 1, Email: "updated@example.com"}

	ctrl := &controllerGeneric[mocks.TestModel, uint]{
		svcT: mockService,
		log:  mockLogger,
	}

	c, w := CreateMockGinContext()

	c.Params = gin.Params{{Key: "id", Value: "1"}}

	mockService.On("Update", c, uint(1), inputModel).Return(gorm.ErrRecordNotFound)

	ctrl.Update(c, inputModel)

	assert.Equal(t, http.StatusNotFound, w.Code)

	expectedBody := `{"err":"no rows found"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}
