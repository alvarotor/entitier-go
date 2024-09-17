package repositories

type IGenericRepo[T any, X string | uint] interface {
	Create(T) (T, error)
	GetAll() ([]*T, error)
	Get(X, string) (*T, error)
	Update(X, T) error
	Delete(X, bool) error
}
