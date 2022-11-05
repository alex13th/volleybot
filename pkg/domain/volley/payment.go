package volley

import (
	"time"
	"volleybot/pkg/domain/person"

	"github.com/google/uuid"
)

type Payment struct {
	Id       uuid.UUID       `json:"id"`
	Person   person.Person   `json:"person"`
	PreCheck PaymentPreCheck `json:"pre_check"`
	Sum      int             `json:"sum"`
}

func (pay Payment) GetSum() int {
	return pay.Sum
}

func (pay Payment) GetId() interface{} {
	return pay.Id
}

func (pay Payment) GetPerson() interface{} {
	return pay.Person
}

type PaymentPreCheck struct {
	Id   string    `json:"id"`
	Time time.Time `json:"time"`
}
