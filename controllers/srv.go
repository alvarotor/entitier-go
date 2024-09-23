package controllers

import (
	"github.com/alvarotor/entitier-go/logger"
	"github.com/alvarotor/entitier-go/repositories"
	"github.com/alvarotor/entitier-go/services"
	"gorm.io/gorm"
)

type controllerGeneric[T any, X string | uint] struct {
	svcT services.GenericService[T, X]
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
