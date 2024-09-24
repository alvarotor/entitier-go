package services

import "github.com/alvarotor/entitier-go/repositories"

type IGenericService[T any, X string | uint] interface {
	repositories.IGenericRepo[T, X]
}
