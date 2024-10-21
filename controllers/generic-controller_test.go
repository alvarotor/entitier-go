package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alvarotor/entitier-go/mocks"
	"github.com/alvarotor/entitier-go/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var ctx = context.Background()

func createMockGinContext() (*gin.Context, *httptest.ResponseRecorder) {
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
		repo: mockService,
		log:  mockLogger,
	}

	c, w := createMockGinContext()

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
		repo: mockService,
		log:  mockLogger,
	}

	c, w := createMockGinContext()

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
		repo: mockService,
		log:  mockLogger,
	}

	c, w := createMockGinContext()

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
		repo: mockService,
		log:  mockLogger,
	}

	result, err := ctrl.Create(ctx, inputModel)

	assert.NoError(t, err)
	assert.Equal(t, createdModel, result)
}

func TestController_Create_Failure(t *testing.T) {
	mockService := new(mocks.IGenericService[mocks.TestModel, uint])
	mockLogger := &mocks.Logger{}

	inputModel := mocks.TestModel{Email: "test@example.com"}

	err := errors.New("create error")
	mockLogger.On("Error", err.Error()).Return(nil)
	mockService.On("Create", ctx, inputModel).Return(inputModel, err)

	ctrl := &controllerGeneric[mocks.TestModel, uint]{
		repo: mockService,
		log:  mockLogger,
	}

	result, errCreate := ctrl.Create(ctx, inputModel)

	assert.Error(t, err)
	assert.Equal(t, inputModel, result)
	assert.Equal(t, err.Error(), errCreate.Error())
}

func TestController_Get_Success(t *testing.T) {
	mockService := new(mocks.IGenericService[mocks.TestModel, uint])
	mockLogger := &mocks.Logger{}

	testModel := &mocks.TestModel{ID: 1, Email: "test1@example.com"}

	ctrl := &controllerGeneric[mocks.TestModel, uint]{
		repo: mockService,
		log:  mockLogger,
	}

	c, w := createMockGinContext()

	c.Set("validatedID", uint(1))

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
				repo: mockService,
				log:  mockLogger,
			}

			c, w := createMockGinContext()

			c.Set("validatedID", uint(1))

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
		repo: mockService,
		log:  mockLogger,
	}

	c, w := createMockGinContext()

	c.Set("validatedID", uint(1))

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
		repo: mockService,
		log:  mockLogger,
	}

	c, w := createMockGinContext()

	c.Set("validatedID", uint(1))

	mockService.On("Delete", c, uint(1), true).Return(models.ErrNotFound)

	ctrl.Delete(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	expectedBody := fmt.Sprintf(`{"err":"%s"}`, models.ErrNotFound.Error())
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestController_Update_Success(t *testing.T) {
	mockService := new(mocks.IGenericService[mocks.TestModel, uint])
	mockLogger := &mocks.Logger{}

	model := mocks.TestModel{ID: 1, Email: "test@example.com"}

	mockService.On("Update", ctx, uint(1), model).Return(nil)

	ctrl := &controllerGeneric[mocks.TestModel, uint]{
		repo: mockService,
		log:  mockLogger,
	}

	status, err := ctrl.Update(ctx, uint(1), model)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
}

func TestController_Update_Failure(t *testing.T) {
	mockService := new(mocks.IGenericService[mocks.TestModel, uint])
	mockLogger := &mocks.Logger{}

	model := mocks.TestModel{ID: 1, Email: "test@example.com"}

	errUpdate := errors.New("update error")
	mockService.On("Update", ctx, uint(1), model).Return(errUpdate)

	ctrl := &controllerGeneric[mocks.TestModel, uint]{
		repo: mockService,
		log:  mockLogger,
	}

	status, err := ctrl.Update(ctx, uint(1), model)

	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Equal(t, errUpdate, err)
}

func TestController_Get_ValidatedIDDoesNotExist(t *testing.T) {
	mockService := new(mocks.IGenericService[mocks.TestModel, uint])
	mockLogger := &mocks.Logger{}

	ctrl := &controllerGeneric[mocks.TestModel, uint]{
		repo: mockService,
		log:  mockLogger,
	}

	c, w := createMockGinContext()

	ctrl.Get(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	expectedBody := fmt.Sprintf(`{"err":"%s"}`, models.ErrMustProvideValidID.Error())
	assert.JSONEq(t, expectedBody, w.Body.String())
}
func TestController_Delete_ValidatedIDDoesNotExist(t *testing.T) {
	mockService := new(mocks.IGenericService[mocks.TestModel, uint])
	mockLogger := &mocks.Logger{}

	ctrl := &controllerGeneric[mocks.TestModel, uint]{
		repo: mockService,
		log:  mockLogger,
	}

	c, w := createMockGinContext()

	ctrl.Delete(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	expectedBody := fmt.Sprintf(`{"err":"%s"}`, models.ErrMustProvideValidID.Error())
	assert.JSONEq(t, expectedBody, w.Body.String())
}
