package models

type PagiationResponse struct {
	Page          int             `json:"page"`
	PageSize      int             `json:"pageSize"`
	TotalElements int             `json:"totalElements"`
	Items         []HotelResponse `json:"items"`
}
