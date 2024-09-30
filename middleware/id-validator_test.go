package middleware

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
)

func setupMockGinContext(method, url string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(method, url, nil)
	return c, w
}

func TestIDValidator_StringID_Success(t *testing.T) {
	mockUtils := new(mocks.MockUtils[string])

	c, w := setupMockGinContext(http.MethodGet, "/test")

	mockUtils.On("GetIDParam", c).Return("abc123")
	mockUtils.On("ConvertToGenericID", "abc123").Return("abc123", nil)

	middleware := IDValidator[string](mockUtils)
	middleware(c)

	validatedID, exists := c.Get("validatedID")
	assert.True(t, exists)
	assert.Equal(t, "abc123", validatedID)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIDValidator_StringID_InvalidID(t *testing.T) {
	mockUtils := new(mocks.MockUtils[string])

	c, w := setupMockGinContext(http.MethodGet, "/test")

	mockUtils.On("GetIDParam", c).Return("invalid-id")
	mockUtils.On("ConvertToGenericID", "invalid-id").Return("", errors.New("invalid string ID"))

	middleware := IDValidator[string](mockUtils)
	middleware(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"err":"invalid string ID"}`, w.Body.String())
}

func TestIDValidator_StringID_MissingIDParam(t *testing.T) {
	mockUtils := new(mocks.MockUtils[string])

	c, w := setupMockGinContext(http.MethodGet, "/test")

	mockUtils.On("GetIDParam", c).Return(nil)

	middleware := IDValidator[string](mockUtils)
	middleware(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, fmt.Sprintf(`{"err":"%s"}`, models.ErrMustProvideValidID.Error()), w.Body.String())
}

func TestIDValidator_UintID_Success(t *testing.T) {
	mockUtils := new(mocks.MockUtils[uint])

	c, w := setupMockGinContext(http.MethodGet, "/test")

	mockUtils.On("GetIDParam", c).Return(uint(123))
	mockUtils.On("ConvertToGenericID", uint(123)).Return(uint(123), nil)

	middleware := IDValidator[uint](mockUtils)
	middleware(c)

	validatedID, exists := c.Get("validatedID")
	assert.True(t, exists)
	assert.Equal(t, uint(123), validatedID)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIDValidator_UintID_InvalidID(t *testing.T) {
	mockUtils := new(mocks.MockUtils[uint])

	c, w := setupMockGinContext(http.MethodGet, "/test")

	mockUtils.On("GetIDParam", c).Return("invalid-id")
	mockUtils.On("ConvertToGenericID", "invalid-id").Return(uint(0), errors.New("invalid uint ID"))

	middleware := IDValidator[uint](mockUtils)
	middleware(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"err":"invalid uint ID"}`, w.Body.String())
}

func TestIDValidator_UintID_MissingIDParam(t *testing.T) {
	mockUtils := new(mocks.MockUtils[uint])

	c, w := setupMockGinContext(http.MethodGet, "/test")

	mockUtils.On("GetIDParam", c).Return(nil)

	middleware := IDValidator[uint](mockUtils)
	middleware(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, fmt.Sprintf(`{"err":"%s"}`, models.ErrMustProvideValidID.Error()), w.Body.String())
}
