package repositories

import (
	"context"
	"errors"
	"reflect"

	"github.com/alvarotor/entitier-go/models"
	"gorm.io/gorm"
)

type genericRepository[T any, X string | uint] struct {
	DB *gorm.DB
}

func NewGenericRepository[T any, X string | uint](db *gorm.DB) IGenericRepo[T, X] {
	return &genericRepository[T, X]{
		DB: db,
	}
}

func (r *genericRepository[T, X]) Create(ctx context.Context, model T) (T, error) {
	// Use reflection to check if model is empty
	if reflect.DeepEqual(model, reflect.Zero(reflect.TypeOf(model)).Interface()) {
		return model, models.ErrModelCannotBeEmpty
	}

	result := r.DB.Create(&model)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return model, models.ErrDuplicatedKeyEmail
		}
		return model, result.Error
	}

	return model, nil
}

func (r genericRepository[T, X]) GetAll(ctx context.Context) ([]*T, error) {
	var items []*T
	result := r.DB.Find(&items)
	if result.Error != nil {
		return items, result.Error
	}
	if len(items) == 0 {
		return items, models.ErrNotFound
	}

	return items, nil
}

func (r *genericRepository[T, X]) Get(ctx context.Context, id X, preload string) (*T, error) {
	var model = new(T)
	result := r.DB
	if len(preload) > 0 {
		result = result.Preload(preload)
	}

	if _, ok := any(id).(string); ok {
		result = result.Where("id = ?", id).First(model)
	} else {
		result = result.First(model, id)
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return model, nil
}

func (r *genericRepository[T, X]) Update(ctx context.Context, id X, amended T) error {
	var existing T
	result := r.DB.First(&existing, "ID = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return models.ErrNotFound
	}
	if result.Error != nil {
		return result.Error
	}

	result = r.DB.Model(&existing).Updates(amended)
	if result.RowsAffected == 0 {
		return models.ErrNotFound
	}
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *genericRepository[T, X]) Delete(ctx context.Context, id X, permanently bool) error {
	t := new(T)
	var deleter *gorm.DB
	if permanently {
		deleter = r.DB.Unscoped()
	} else {
		deleter = r.DB
	}

	if _, ok := any(id).(string); ok {
		deleter = deleter.Where("id = ?", id).Delete(t)
	} else {
		deleter = deleter.Delete(t, id)
	}

	if deleter.Error != nil {
		if errors.Is(deleter.Error, gorm.ErrRecordNotFound) {
			return models.ErrNotFound
		}
		return deleter.Error
	}

	if deleter.RowsAffected == 0 {
		return models.ErrNotFound
	}

	return nil
}
