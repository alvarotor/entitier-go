package controllers

import (
	"context"
	"errors"
	"net/http"

	"github.com/alvarotor/entitier-go/logger"
	"github.com/alvarotor/entitier-go/models"
	"github.com/alvarotor/entitier-go/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type controllerGeneric[T any, X string | uint] struct {
	repo repository.IGenericRepo[T, X]
	log  logger.Logger
}

func NewGenericController[T any, X string | uint](log logger.Logger, db *gorm.DB) IControllerGeneric[T, X] {
	repo := repository.NewGenericRepository[T, X](
		db,
	)
	return &controllerGeneric[T, X]{
		repo: repo,
		log:  log,
	}
}

func (u *controllerGeneric[T, X]) Create(ctx context.Context, model T) (T, error) {
	m, err := u.repo.Create(ctx, model)
	if err != nil {
		u.log.Error("create", err.Error())
		return m, err
	}

	return m, nil
}

func (u *controllerGeneric[T, X]) Get(c *gin.Context) {
	id, exists := c.Get("validatedID")
	if !exists {
		handleError(c, u.log, "get", models.ErrMustProvideValidID, http.StatusBadRequest)
		return
	}

	preloadArgInterface, exists := c.Get("preloadArg")
	var preloadArg string
	if exists && preloadArgInterface != nil {
		preloadArg = preloadArgInterface.(string)
	} else {
		preloadArg = ""
	}

	p, err := u.repo.Get(c, id.(X), preloadArg)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			handleError(c, u.log, "get", models.ErrNotFound, http.StatusNotFound)
		} else {
			handleError(c, u.log, "get", err, http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": p})
}

func (u *controllerGeneric[T, X]) GetAll(c *gin.Context) {
	ps, err := u.repo.GetAll(c)
	if errors.Is(err, models.ErrNotFound) {
		handleError(c, u.log, "getall", err, http.StatusNotFound)
		return
	}
	if err != nil {
		handleError(c, u.log, "getall", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"all": ps})
}

func (u *controllerGeneric[T, X]) Delete(c *gin.Context) {
	id, exists := c.Get("validatedID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"err": models.ErrMustProvideValidID.Error()})
		return
	}

	err := u.repo.Delete(c, id.(X), true)
	if err != nil {
		handleError(c, u.log, "delete", err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (u *controllerGeneric[T, X]) Update(ctx context.Context, id X, model T) (int, error) {
	err := u.repo.Update(ctx, id, model)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func handleError(c *gin.Context, log logger.Logger, id string, err error, statusCode int) {
	log.Error(id, err.Error())
	c.JSON(statusCode, gin.H{"err": err.Error()})
}
