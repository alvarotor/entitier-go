package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (u *controllerGeneric[T, X]) GetAll(c *gin.Context) {
	ps, err := u.svcT.GetAll()
	if err != nil {
		u.log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"all": ps})
}
