package controllers

import (
	"github.com/alvarotor/entitier-go/logger"
	"github.com/alvarotor/entitier-go/repositories"
	"github.com/alvarotor/entitier-go/services"
	"gorm.io/gorm"
)

type controllerGen[T any, X string | uint] struct {
	svcT services.GenericService[T, X]
	log  logger.Logger
}

func NewGenericControllerSrv[T any, X string | uint](log logger.Logger, db *gorm.DB) IControllerGen[T, X] {
	repo := repositories.NewGenericRepository[T, X](
		db,
	)
	svcGen := services.NewGenericService(
		repo,
	)
	return &controllerGen[T, X]{
		svcT: svcGen,
		log:  log,
	}
}
