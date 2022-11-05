package order

import (
	uuid "github.com/google/uuid"
)

type OrderRepository interface {
	Get(uuid.UUID) (Order, error)
	Add(Order) (Order, error)
	Update(Order) error
}

type PaymentRepository interface {
	Get(Payment) (Payment, error)
	Add(Payment) (Order, error)
	Update(Payment) error
}
