package order

import (
	"time"
	"volleybot/pkg/domain/person"

	"github.com/google/uuid"
)

type Payment struct {
	Id     uuid.UUID     `json:"id"`
	Person person.Person `json:"person"`
	Sum    int           `json:"sum"`
}

type TelegramPay struct {
	Payment
	PreCheck PaymentPreCheck `json:"pre_check"`
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
