package utils

import (
	"strconv"

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
