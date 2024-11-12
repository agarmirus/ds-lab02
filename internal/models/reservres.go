package models

type ReservationResponse struct {
	ReservationUid string      `json:"reservationUid"`
	Hotel          HotelInfo   `json:hotel`
	StartDate      string      `json:"startDate"`
	EndDate        string      `json:"endDate"`
	Status         string      `json:"status"`
	Payment        PaymentInfo `json:"payment"`
}
