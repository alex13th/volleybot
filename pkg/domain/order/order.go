package order

import (
	"time"
	"volleybot/pkg/domain/person"

	"github.com/google/uuid"
)

type Order struct {
	Id       uuid.UUID     `json:"id"`
	Date     time.Time     `json:"date"`
	Person   person.Person `json:"person"`
	Sum      int           `json:"sum"`
	Payments []Payment     `json:"payments"`
}

type Payment interface {
	GetId() interface{}
	GetPerson() person.Person
	GetSum() int
}
