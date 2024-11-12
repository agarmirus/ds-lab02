package models

type Payment struct {
	Id     int    `json:"id"`
	Uid    string `json:"paymentUid"`
	Status string `json:"status"`
	Price  int    `json:"price"`
}
