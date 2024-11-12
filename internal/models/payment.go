package models

import "github.com/google/uuid"

type Payment struct {
	id     int
	uid    uuid.UUID
	status string
	price  int
}

func (payment *Payment) SetId(id int) {
	payment.id = id
}

func (payment *Payment) SetUid(uid uuid.UUID) {
	payment.uid = uid
}

func (payment *Payment) SetStatus(status string) {
	payment.status = status
}

func (payment *Payment) SetPrice(price int) {
	payment.price = price
}

func (payment Payment) GetId() int {
	return payment.id
}

func (payment Payment) GetUid() uuid.UUID {
	return payment.uid
}

func (payment Payment) GetStatus() string {
	return payment.status
}

func (payment Payment) GetPrice() int {
	return payment.price
}
