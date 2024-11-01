package database

import "container/list"

type IDAO[T any] interface {
	Create(T) (T, error)

	GetById(T) (T, error)
	GetByAttribute(string, string) (list.List, error)

	Update(T) (T, error)

	Delete(T) error
}
