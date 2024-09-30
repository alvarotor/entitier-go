package middleware

import (
	"net/http"

	"github.com/alvarotor/entitier-go/models"
	"github.com/gin-gonic/gin"
)

type UtilsInterface[X string | uint] interface {
	GetIDParam(c *gin.Context) interface{}
	ConvertToGenericID(idInterface interface{}) (X, error)
}

func IDValidator[X string | uint](utils UtilsInterface[X]) gin.HandlerFunc {
	return func(c *gin.Context) {
		idInterface := utils.GetIDParam(c)
		if idInterface == nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": models.ErrMustProvideValidID.Error()})
			c.Abort()
			return
		}

		id, err := utils.ConvertToGenericID(idInterface)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			c.Abort()
			return
		}

		c.Set("validatedID", id)
		c.Next()
	}
}
