package models

import (
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	id        int
	uid       uuid.UUID
	username  string
	payment   Payment
	hotel     Hotel
	status    string
	startDate time.Time
	endDate   time.Time
}

func (reservation *Reservation) SetId(id int) {
	reservation.id = id
}

func (reservation *Reservation) SetUid(uid uuid.UUID) {
	reservation.uid = uid
}

func (reservation *Reservation) SetUsername(username string) {
	reservation.username = username
}

func (reservation *Reservation) SetPaymentUid(payment Payment) {
	reservation.payment = payment
}

func (reservation *Reservation) SetHotelId(hotel Hotel) {
	reservation.hotel = hotel
}

func (reservation *Reservation) SetStatus(status string) {
	reservation.status = status
}

func (reservation *Reservation) SetStartDate(startDate time.Time) {
	reservation.startDate = startDate
}

func (reservation *Reservation) SetEndDate(endDate time.Time) {
	reservation.endDate = endDate
}

func (reservation *Reservation) GetId() int {
	return reservation.id
}

func (reservation *Reservation) GetUid() uuid.UUID {
	return reservation.uid
}

func (reservation *Reservation) GetUsername() string {
	return reservation.username
}

func (reservation *Reservation) GetPaymentUid() Payment {
	return reservation.payment
}

func (reservation *Reservation) GetHotelId() Hotel {
	return reservation.hotel
}

func (reservation *Reservation) GetStatus() string {
	return reservation.status
}

func (reservation *Reservation) GetStartDate() time.Time {
	return reservation.startDate
}

func (reservation *Reservation) GetEndDate() time.Time {
	return reservation.endDate
}
