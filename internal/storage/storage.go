package storage

import (
	"errors"

	"github.com/voznikaetnepriyazn/Good-service/internal/models/good"

	"github.com/google/uuid"
)

var (
	ErrUrlNotFound = errors.New("url not found")
	ErrUrlExist    = errors.New("url exist")
)

type GoodService interface {
	AddURL(good good.Good) (uuid.UUID, error)

	DeleteURL(id uuid.UUID) error

	GetAllURL() ([]good.Good, error)

	GetByIdURL(id uuid.UUID) (uuid.UUID, error)

	UpdateURL(good good.Good) error

	GetListOfGoodsByBrand(id uuid.UUID) ([]good.Good, error)

	GetListOfGoodsByType(id uuid.UUID) ([]good.Good, error)

	IsAvaliableForOrder(id uuid.UUID) (bool, error)

	RestOfGood(id uuid.UUID) (int, error)
}
