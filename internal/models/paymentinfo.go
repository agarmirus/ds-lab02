package models

type PaymentInfo struct {
	Status string `json:"status"`
	Price  int    `json:"price"`
}
