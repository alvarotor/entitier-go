package repositories

import (
	"errors"

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

func (r *genericRepository[T, X]) Create(model T) (T, error) {
	result := r.DB.Create(&model)

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return model, models.ErrDuplicatedKeyEmail
	}
	if result.Error != nil {
		return model, result.Error
	}
	if result.RowsAffected == 0 {
		return model, models.ErrNotFound
	}

	return model, nil
}

func (r *genericRepository[T, X]) GetAll() ([]*T, error) {
	model := []*T{}
	result := r.DB.Find(&model)
	if result.RowsAffected == 0 {
		return model, models.ErrNotFound
	}
	if result.Error != nil {
		return model, result.Error
	}

	return model, nil
}

func (r *genericRepository[T, X]) Get(id X, preload string) (*T, error) {
	var model = new(T)
	result := r.DB
	if len(preload) > 0 {
		result = result.Preload(preload)
	}
	result = result.First(&model, "ID = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return model, gorm.ErrRecordNotFound
	}
	if result.Error != nil {
		return model, result.Error
	}

	return model, nil
}

func (r *genericRepository[T, X]) Update(id X, amended T) error {
	result := r.DB.Save(&amended)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return gorm.ErrRecordNotFound
	}
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *genericRepository[T, X]) Delete(id X, permanently bool) error {
	t := new(T)
	var deleter *gorm.DB
	if permanently {
		deleter = r.DB.Unscoped()
	} else {
		deleter = r.DB
	}

	switch any(id).(type) {
	case string:
		deleter = deleter.Where("id = ?", id).Delete(t)
	default:
		deleter = deleter.Delete(t, id)
	}
	if errors.Is(deleter.Error, gorm.ErrRecordNotFound) {
		return gorm.ErrRecordNotFound
	}
	if deleter.Error != nil {
		return deleter.Error
	}

	return nil
}
