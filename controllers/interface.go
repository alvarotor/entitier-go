package controllers

import "github.com/gin-gonic/gin"

type IControllerGeneric[T any, X string | uint] interface {
	GetAll(c *gin.Context)
	Create(T) (T, error)
	Get(c *gin.Context)
}
