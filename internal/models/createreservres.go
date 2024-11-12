package models

type CreateReservationResponse struct {
	ReservationUid string      `json:"reservationUid"`
	HotelUid       string      `json:"hotelUid"`
	StartDate      string      `json:"startDate"`
	EndDate        string      `json:"endDate"`
	Discount       int         `json:"discount"`
	Status         string      `json:"status"`
	Payment        PaymentInfo `json:"payment"`
}
