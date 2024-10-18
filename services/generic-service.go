package services

import (
	"context"

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

func (r *GenericService[T, X]) GetAll(ctx context.Context) ([]*T, error) {
	model, err := r.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (r *GenericService[T, X]) Get(ctx context.Context, ID X, preload string) (*T, error) {
	model, err := r.repo.Get(ctx, ID, preload)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (r *GenericService[T, X]) Create(ctx context.Context, data T) (T, error) {
	created, err := r.repo.Create(ctx, data)
	if err != nil {
		return data, err
	}

	return created, nil
}

func (r *GenericService[T, X]) Delete(ctx context.Context, ID X, permanently bool) error {
	err := r.repo.Delete(ctx, ID, permanently)
	if err != nil {
		return err
	}

	return nil
}

func (r *GenericService[T, X]) Update(ctx context.Context, ID X, amended T) error {
	err := r.repo.Update(ctx, ID, amended)
	if err != nil {
		return err
	}

	return nil
}
