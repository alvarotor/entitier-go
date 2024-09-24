package services

import (
	"github.com/alvarotor/entitier-go/repositories"
)

type GenericService[T any, X string | uint] struct {
	repo repositories.IGenericRepo[T, X]
}

func NewGenericService[T any, X string | uint](
	repo repositories.IGenericRepo[T, X],
) IGenericService[T, X] {
	return &GenericService[T, X]{
		repo: repo,
	}
}

func (r *GenericService[T, X]) GetAll() ([]*T, error) {
	model, err := r.repo.GetAll()
	if err != nil {
		return nil, err
	}

	if len(model) == 0 {
		return nil, err
	}

	return model, nil
}

func (r *GenericService[T, X]) Get(ID X, preload string) (*T, error) {
	model, err := r.repo.Get(ID, preload)
	if err != nil {
		return model, err
	}

	return model, nil
}

func (r *GenericService[T, X]) Create(data T) (T, error) {
	created, err := r.repo.Create(data)
	if err != nil {
		return data, err
	}

	return created, nil
}

func (r *GenericService[T, X]) Delete(ID X, permanently bool) error {
	err := r.repo.Delete(ID, permanently)
	if err != nil {
		return err
	}

	return nil
}

func (r *GenericService[T, X]) Update(ID X, amended T) error {
	err := r.repo.Update(ID, amended)
	if err != nil {
		return err
	}

	return nil
}
