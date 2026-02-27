package typee

import (
	"github.com/google/uuid"
	"github.com/voznikaetnepriyazn/Good-service/internal/models/good"
)

type Typee struct {
	Id         uuid.UUID
	Name       string
	ListOfGood []*good.Good
}
