package repositories

import "context"

type IGenericRepo[T any, X string | uint] interface {
	Create(context.Context, T) (T, error)
	GetAll(context.Context) ([]*T, error)
	Get(context.Context, X, string) (*T, error)
	Update(context.Context, X, T) error
	Delete(context.Context, X, bool) error
}
