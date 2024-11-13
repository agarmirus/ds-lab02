package models

import (
	"time"
)

type Reservation struct {
	Id         int       `json:"id"`
	Uid        string    `json:"reservationUid"`
	Username   string    `json:"username"`
	PaymentUid string    `json:"paymentUid"`
	HotelId    int       `json:"hotelId"`
	Status     string    `json:"status"`
	StartDate  time.Time `json:"startDate"`
	EndDate    time.Time `json:"endDate"`
}
