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

func ConvertToGenericID[X string | uint](id interface{}) (X, error) {
	var zeroX X

	switch v := id.(type) {
	case string:
		if _, isString := any(zeroX).(string); isString {
			return any(v).(X), nil
		}
	case uint:
		if _, isUint := any(zeroX).(uint); isUint {
			return any(v).(X), nil
		}
	}
	return zeroX, models.ErrIDTypeMismatch
}
