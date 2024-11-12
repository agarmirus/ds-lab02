package models

type CreateReservationRequest struct {
	HotelUid  string `json:"hotelUid"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}
