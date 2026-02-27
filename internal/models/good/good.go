package good

import (
	"github.com/google/uuid"
)

type Good struct {
	Id      uuid.UUID
	name    string
	TypeId  uuid.UUID
	BrandId uuid.UUID
	Rest    int16
}
