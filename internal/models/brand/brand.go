package brand

import (
	"github.com/google/uuid"
	"github.com/voznikaetnepriyazn/Good-service/internal/models/good"
)

type Brand struct {
	Id         uuid.UUID
	Name       string
	ListOfGood []*good.Good
}
