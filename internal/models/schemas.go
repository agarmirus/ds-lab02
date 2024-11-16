package models

import (
	"errors"
	"time"

	"github.com/agarmirus/ds-lab02/internal/serverrors"
	"github.com/google/uuid"
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

type Payment struct {
	Id     int    `json:"id"`
	Uid    string `json:"paymentUid"`
	Status string `json:"status"`
	Price  int    `json:"price"`
}

type Loyalty struct {
	Id               int    `json:"id"`
	Username         string `json:"username"`
	ReservationCount int    `json:"reservationCount"`
	Status           string `json:"status"`
	Discount         int    `json:"discount"`
}

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

type PaymentInfo struct {
	Status string `json:"status"`
	Price  int    `json:"price"`
}

type LoyaltyInfoResponse struct {
	Status           string `json:"status"`
	Discount         int    `json:"discount"`
	ReservationCount int    `json:"reservationCount"`
}

type HotelResponse struct {
	HotelUid string `json:"hotelUid"`
	Name     string `json:"name"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Address  string `json:"address"`
	Stars    int    `json:"stars"`
	Price    int    `json:"price"`
}

type HotelInfo struct {
	HotelUid    string `json:"hotelUid"`
	Name        string `json:"name"`
	FullAddress string `json:"fullAddress"`
	Stars       int    `json:"stars"`
}

type ReservationResponse struct {
	ReservationUid string      `json:"reservationUid"`
	Hotel          HotelInfo   `json:"hotel"`
	StartDate      string      `json:"startDate"`
	EndDate        string      `json:"endDate"`
	Status         string      `json:"status"`
	Payment        PaymentInfo `json:"payment"`
}

type UserInfoResponse struct {
	Reservations []ReservationResponse `json:"reservations"`
	Loyalty      LoyaltyInfoResponse   `json:"loyalty"`
}

type CreateReservationRequest struct {
	HotelUid  string `json:"hotelUid"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type CreateReservationResponse struct {
	ReservationUid string      `json:"reservationUid"`
	HotelUid       string      `json:"hotelUid"`
	StartDate      string      `json:"startDate"`
	EndDate        string      `json:"endDate"`
	Discount       int         `json:"discount"`
	Status         string      `json:"status"`
	Payment        PaymentInfo `json:"payment"`
}

type PagiationResponse struct {
	Page          int             `json:"page"`
	PageSize      int             `json:"pageSize"`
	TotalElements int             `json:"totalElements"`
	Items         []HotelResponse `json:"items"`
}

type ErrorDiscription struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ValidationErrorResponse struct {
	Message string             `json:"message"`
	Errors  []ErrorDiscription `json:"errors"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func HotelToHotelReponse(
	hotelRes *HotelResponse,
	hotel *Hotel,
) {
	hotelRes.HotelUid = hotel.Uid
	hotelRes.Name = hotel.Name
	hotelRes.Country = hotel.Country
	hotelRes.City = hotel.City
	hotelRes.Address = hotel.Address
	hotelRes.Stars = hotel.Stars
	hotelRes.Price = hotel.Price
}

func HotelsSliceToPagRes(
	pagRes *PagiationResponse,
	hotelsSlice []Hotel,
	page int,
	pageSize int,
) {
	pagRes.Page = page
	pagRes.PageSize = pageSize
	pagRes.TotalElements = len(hotelsSlice)

	for i := 0; i < pagRes.TotalElements; i++ {
		var hotelRes HotelResponse
		HotelToHotelReponse(&hotelRes, &hotelsSlice[i])

		pagRes.Items = append(pagRes.Items, hotelRes)
	}
}

func ReservsSliceToUserInfoRes(
	userInfoRes *UserInfoResponse,
	userReservsSlice []Reservation,
	hotelsMap map[int]Hotel,
	paymentsMap map[string]Payment,
	loyalty *Loyalty,
) {
}

func ReservToReservRes(
	reservRes *ReservationResponse,
	reservation *Reservation,
	hotel *Hotel,
	payment *Payment,
) {
}

func ReservsSliceToReservRes(
	reservsResSlice []ReservationResponse,
	userReservsSlice []Reservation,
	hotelsMap map[int]Hotel,
	paymentsMap map[string]Payment,
) {
}

func ReservToCrReservRes(
	crReservRes *CreateReservationResponse,
	reservation *Reservation,
	payment *Payment,
	loyalty *Loyalty,
) {
}

func LoyaltyToLoyaltyInfoRes(
	loyatlyInfoRes *LoyaltyInfoResponse,
	loyalty *Loyalty,
) {
}

func ValidateCrReservReq(
	createReservReq *CreateReservationRequest,
) (validErrRes ValidationErrorResponse, err error) {
	if uuid.Validate(createReservReq.HotelUid) != nil {
		validErrRes.Errors = append(validErrRes.Errors, ErrorDiscription{Field: `hotelUid`, Error: `invalid uid`})
	}

	startDate, err := time.Parse(`%F`, createReservReq.StartDate)

	if err != nil {
		validErrRes.Errors = append(validErrRes.Errors, ErrorDiscription{Field: `startDate`, Error: `invalid date format`})
	}

	endDate, err := time.Parse(`%F`, createReservReq.EndDate)

	if err != nil {
		validErrRes.Errors = append(validErrRes.Errors, ErrorDiscription{Field: `endDate`, Error: `invalid date format`})
	} else if startDate.Unix() > endDate.Unix() {
		validErrRes.Errors = append(validErrRes.Errors, ErrorDiscription{Field: `startDate`, Error: `invalid date period`})
		err = errors.New(serverrors.ErrInvalidReservDates)
	}

	if err != nil {
		validErrRes.Message = `invalid reservation request data`
	}

	return validErrRes, err
}
