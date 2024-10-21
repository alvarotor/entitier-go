package controllers

import (
	"context"

	"github.com/gin-gonic/gin"
)

type IControllerGeneric[T any, X string | uint] interface {
	GetAll(*gin.Context)
	Create(context.Context, T) (T, error)
	Get(*gin.Context)
	Delete(*gin.Context)
	Update(context.Context, X, T) (int, error)
}
