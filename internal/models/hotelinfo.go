package models

type HotelInfo struct {
	HotelUid    string `json:"hotelUid"`
	Name        string `json:"name"`
	FullAddress string `json:fullAddress`
	Stars       int    `json:stars`
}
