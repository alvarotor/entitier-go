package models

import "errors"

var (
	ErrNotFound           = errors.New("no rows found")
	ErrDuplicatedKeyEmail = errors.New("duplicated key. Email already exists")
	ErrModelCannotBeEmpty = errors.New("model cannot be empty")
	ErrMustProvideValidID = errors.New("must provide valid id")
)
