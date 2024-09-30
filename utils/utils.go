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
	if strVal, ok := idInterface.(string); ok {
		var zeroX X
		if _, isString := any(zeroX).(string); isString {
			return any(strVal).(X), nil
		}
	}

	if uintVal, ok := idInterface.(uint); ok {
		var zeroX X
		if _, isUint := any(zeroX).(uint); isUint {
			return any(uintVal).(X), nil
		}
	}

	return *new(X), models.ErrIDTypeMismatch
}
