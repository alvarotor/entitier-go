package utils_test

import (
	"testing"

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
