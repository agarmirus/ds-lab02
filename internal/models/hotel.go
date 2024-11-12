package models

type Hotel struct {
	Id      int    `json:"id"`
	Uid     string `json:"hotelUid"`
	Name    string `json:"name"`
	Country string `json:"country"`
	City    string `json:"city"`
	Address string `json:"address"`
	Stars   int    `json:"stars"`
	Price   int    `json:"price"`
}
