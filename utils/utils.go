package utils

import (
	"strconv"

	"github.com/alvarotor/entitier-go/models"
	"github.com/gin-gonic/gin"
)

func GetIDParam(c *gin.Context) interface{} {
	idStr := c.Param("id")
	if idStr == "" {
		return nil
	}

	if idUint, err := strconv.ParseUint(idStr, 10, 64); err == nil {
		return uint(idUint)
	}

	return idStr
}

func ConvertToGenericID[X string | uint](idInterface interface{}) (X, error) {
	// Check if idInterface is a string
	if strVal, ok := idInterface.(string); ok {
		var zeroX X
		// Check if X is string and assign
		if _, isString := any(zeroX).(string); isString {
			return any(strVal).(X), nil // Type assertion to cast string to X
		}
	}

	// Check if idInterface is a uint
	if uintVal, ok := idInterface.(uint); ok {
		var zeroX X
		// Check if X is uint and assign
		if _, isUint := any(zeroX).(uint); isUint {
			return any(uintVal).(X), nil // Type assertion to cast uint to X
		}
	}

	return *new(X), models.ErrIDTypeMismatch
}
