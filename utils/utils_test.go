package utils_test

import (
	"testing"

	"github.com/alvarotor/entitier-go/models"
	"github.com/alvarotor/entitier-go/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Utility to create a mock Gin context with given parameters
func createGinContextWithParam(key, value string) *gin.Context {
	c, _ := gin.CreateTestContext(nil)
	if key != "" && value != "" {
		c.Params = gin.Params{
			{Key: key, Value: value},
		}
	}
	return c
}

func TestGetIDParam_MissingIDParam(t *testing.T) {
	// Test when the "id" param is missing
	c := createGinContextWithParam("", "") // No param set

	result := utils.GetIDParam(c)

	// Expecting nil when there's no "id" param
	assert.Nil(t, result, "Expected result to be nil when 'id' param is missing")
}

func TestGetIDParam_ValidIDParam(t *testing.T) {
	// Test when the "id" param is a valid unsigned integer
	c := createGinContextWithParam("id", "123")

	result := utils.GetIDParam(c)

	// Expecting uint(123)
	assert.IsType(t, uint(0), result, "Expected result to be of type uint")
	assert.Equal(t, uint(123), result, "Expected result to be 123")
}

func TestGetIDParam_InvalidIDParam(t *testing.T) {
	// Test when the "id" param is an invalid string (not a number)
	c := createGinContextWithParam("id", "abc123")

	result := utils.GetIDParam(c)

	// Expecting the raw string "abc123"
	assert.IsType(t, "", result, "Expected result to be of type string")
	assert.Equal(t, "abc123", result, "Expected result to be 'abc123'")
}

func TestGetIDParam_EmptyIDParam(t *testing.T) {
	// Test when the "id" param is an empty string
	c := createGinContextWithParam("id", "")

	result := utils.GetIDParam(c)

	// Expecting nil since the param is empty
	assert.Nil(t, result, "Expected result to be nil when 'id' param is an empty string")
}

func TestGetIDParam_OverflowIDParam(t *testing.T) {
	// Test when the "id" param is a number that overflows uint64
	// This is useful to ensure that large numbers are handled correctly
	c := createGinContextWithParam("id", "18446744073709551616") // One more than uint64 max value

	result := utils.GetIDParam(c)

	// Since it overflows, expecting the raw string
	assert.IsType(t, "", result, "Expected result to be of type string for overflow")
	assert.Equal(t, "18446744073709551616", result, "Expected result to be the string representing the large number")
}

// Tests for ConvertToGenericID

func TestConvertToGenericID(t *testing.T) {
	t.Run("ValidStringID", func(t *testing.T) {
		idInterface := "12345"
		expectedID := "12345"

		id, err := utils.ConvertToGenericID[string](idInterface)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, id)
	})

	t.Run("ValidUintID", func(t *testing.T) {
		idInterface := uint(12345)
		expectedID := uint(12345)

		id, err := utils.ConvertToGenericID[uint](idInterface)

		assert.NoError(t, err)
		assert.Equal(t, expectedID, id)
	})

	t.Run("InvalidTypeMismatchString", func(t *testing.T) {
		idInterface := uint(12345)

		_, err := utils.ConvertToGenericID[string](idInterface)

		assert.Error(t, err)
		assert.Equal(t, models.ErrIDTypeMismatch, err)
	})

	t.Run("InvalidTypeMismatchUint", func(t *testing.T) {
		idInterface := "12345"

		_, err := utils.ConvertToGenericID[uint](idInterface)

		assert.Error(t, err)
		assert.Equal(t, models.ErrIDTypeMismatch, err)
	})

	t.Run("NilID", func(t *testing.T) {
		var idInterface interface{}

		id, err := utils.ConvertToGenericID[string](idInterface)

		assert.Error(t, err)
		assert.Equal(t, models.ErrIDTypeMismatch, err)
		assert.Equal(t, "", id) // since X is string, expect empty string as the default zero value
	})

	t.Run("EmptyStringID", func(t *testing.T) {
		idInterface := ""

		id, err := utils.ConvertToGenericID[string](idInterface)

		assert.NoError(t, err)
		assert.Equal(t, "", id) // Empty string is still valid as a string ID
	})

	t.Run("ZeroUintID", func(t *testing.T) {
		idInterface := uint(0)

		id, err := utils.ConvertToGenericID[uint](idInterface)

		assert.NoError(t, err)
		assert.Equal(t, uint(0), id) // Zero uint is still valid as a uint ID
	})

	t.Run("UnsupportedTypeFloat", func(t *testing.T) {
		idInterface := 123.45

		_, err := utils.ConvertToGenericID[string](idInterface)

		assert.Error(t, err)
		assert.Equal(t, models.ErrIDTypeMismatch, err)
	})

	t.Run("UnsupportedTypeBool", func(t *testing.T) {
		idInterface := true

		_, err := utils.ConvertToGenericID[uint](idInterface)

		assert.Error(t, err)
		assert.Equal(t, models.ErrIDTypeMismatch, err)
	})
}
