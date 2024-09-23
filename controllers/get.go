package controllers

import (
	"net/http"

	"github.com/alvarotor/entitier-go/utils"
	"github.com/gin-gonic/gin"
)

func (u *controllerGeneric[T, X]) Get(c *gin.Context) {
	idInterface := utils.GetIDParam(c)
	if idInterface == nil {
		err := "must provide valid id"
		u.log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}

	var id X
	var ok bool

	id, ok = idInterface.(X)
	if !ok {
		err := "id type mismatch"
		u.log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}

	p, err := u.svcT.Get(id, "User")
	if err != nil {
		u.log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": p})
}
