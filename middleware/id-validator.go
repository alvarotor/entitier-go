package middleware

import (
	"net/http"
	"strconv"

	"github.com/alvarotor/entitier-go/models"
	"github.com/gin-gonic/gin"
)

func IDValidator[X string | uint]() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		if idStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"err": models.ErrMustProvideValidID.Error()})
			c.Abort()
			return
		}

		idInterface := getIDParam(c)
		id, err := convertToGenericID[X](idInterface)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			c.Abort()
			return
		}

		c.Set("validatedID", id)
		c.Next()
	}
}

func getIDParam(c *gin.Context) interface{} {
	idStr := c.Param("id")
	if idUint, err := strconv.ParseUint(idStr, 10, 64); err == nil {
		return uint(idUint)
	}
	return idStr
}

func convertToGenericID[X string | uint](id interface{}) (X, error) {
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
