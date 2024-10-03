package controllers

import (
	"context"

	"github.com/gin-gonic/gin"
)

type IControllerGeneric[T any, X string | uint] interface {
	GetAll(c *gin.Context)
	Create(context.Context, T) (T, error)
	Get(c *gin.Context)
	Delete(c *gin.Context)
	Update(ctx context.Context, id X, model T) (int, error)
}
