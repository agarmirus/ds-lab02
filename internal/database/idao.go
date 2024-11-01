package database

import "container/list"

type IDAO[T any] interface {
	SetConnectionString(string)

	Create(*T) error

	Get() (list.List, error)
	GetById(*T) (T, error)
	GetByAttribute(string, string) (list.List, error)

	Update(*T) (T, error)

	Delete(*T) error
}
