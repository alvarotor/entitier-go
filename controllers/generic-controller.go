package controllers

import (
	"errors"
	"net/http"

	"github.com/alvarotor/entitier-go/logger"
	"github.com/alvarotor/entitier-go/models"
	"github.com/alvarotor/entitier-go/repositories"
	"github.com/alvarotor/entitier-go/services"
	"github.com/alvarotor/entitier-go/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IControllerGeneric[T any, X string | uint] interface {
	GetAll(c *gin.Context)
	Create(T) (T, error)
	Get(c *gin.Context)
	Delete(c *gin.Context)
	Update(c *gin.Context, model T)
}

type controllerGeneric[T any, X string | uint] struct {
	svcT services.IGenericService[T, X]
	log  logger.Logger
}

func NewGenericController[T any, X string | uint](log logger.Logger, db *gorm.DB) IControllerGeneric[T, X] {
	repo := repositories.NewGenericRepository[T, X](
		db,
	)
	svcGen := services.NewGenericService(
		repo,
	)
	return &controllerGeneric[T, X]{
		svcT: svcGen,
		log:  log,
	}
}

func (u *controllerGeneric[T, X]) Create(model T) (T, error) {
	m, err := u.svcT.Create(model)
	if err != nil {
		u.log.Error(err.Error())
		return m, err
	}

	return m, nil
}

func (u *controllerGeneric[T, X]) Get(c *gin.Context) {
	idInterface := utils.GetIDParam(c)
	if idInterface == nil {
		err := models.ErrMustProvideValidID
		u.log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	id, err := utils.ConvertToGenericID[X](idInterface)
	if err != nil {
		u.log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	p, err := u.svcT.Get(id, "User")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"err": models.ErrNotFound.Error()}) // Ensure error is a string here
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()}) // Ensure proper string conversion for other errors
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"item": p})
}

func (u *controllerGeneric[T, X]) GetAll(c *gin.Context) {
	ps, err := u.svcT.GetAll()
	if errors.Is(err, models.ErrNotFound) {
		u.log.Error(err.Error())
		c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
		return
	}
	if err != nil {
		u.log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"all": ps})
}

func (u *controllerGeneric[T, X]) Delete(c *gin.Context) {
	idInterface := utils.GetIDParam(c)
	if idInterface == nil {
		err := models.ErrMustProvideValidID
		u.log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	id, err := utils.ConvertToGenericID[X](idInterface)
	if err != nil {
		u.log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = u.svcT.Delete(id, true)
	if err != nil {
		u.log.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (u *controllerGeneric[T, X]) Update(c *gin.Context, model T) {
	idInterface := utils.GetIDParam(c)
	if idInterface == nil {
		err := models.ErrMustProvideValidID
		u.log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	id, err := utils.ConvertToGenericID[X](idInterface)
	if err != nil {
		u.log.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = u.svcT.Update(id, model)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"err": models.ErrNotFound.Error()}) // Ensure error is a string here
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()}) // Ensure proper string conversion for other errors
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}
