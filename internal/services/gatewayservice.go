package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/agarmirus/ds-lab02/internal/models"
	"github.com/agarmirus/ds-lab02/internal/serverrors"
	"github.com/google/uuid"
)

type GatewayService struct {
	reservServiceHost string
	reservServicePort int

	paymentServiceHost string
	paymentServicePort int

	loyaltyServiceHost string
	loyaltyServicePost int
}

func NewGatewayService(
	reservServiceHost string,
	reservServicePort int,
	paymentServiceHost string,
	paymentServicePort int,
	loyaltyServiceHost string,
	loyaltyServicePost int,
) IGatewayService {
	return &GatewayService{
		reservServiceHost,
		reservServicePort,
		paymentServiceHost,
		paymentServicePort,
		loyaltyServiceHost,
		loyaltyServicePost,
	}
}

func (service *GatewayService) performAllHotelsGetRequest() (hotelsSlice []models.Hotel, err error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"http://%s:%d/api/v1/hotels",
			service.reservServiceHost,
			service.reservServicePort,
		),
		nil,
	)

	if err != nil {
		return hotelsSlice, errors.New(serverrors.ErrNewRequestForming)
	}

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return hotelsSlice, errors.New(serverrors.ErrRequestSend)
	}

	var resBody []byte
	_, err = res.Body.Read(resBody)

	if err != nil {
		return hotelsSlice, errors.New(serverrors.ErrResponseRead)
	}

	err = json.Unmarshal(resBody, &hotelsSlice)

	if err != nil {
		return hotelsSlice, errors.New(serverrors.ErrResponseParse)
	}

	return hotelsSlice, nil
}

func (service *GatewayService) performUserReservsGetRequest(
	username string,
) (userReservsSlice []models.Reservation, err error) {
}

func (service *GatewayService) performReservsHotelsGetRequest(
	userReservsSlice []models.Reservation,
) (hotelsMap map[int]models.Hotel, err error) {
}

func (service *GatewayService) performPaymentByUidGet(
	paymentUid string,
) (payment models.Payment, err error) {
}

func (service *GatewayService) performReservsPaymentsGetRequest(
	userReservsSlice []models.Reservation,
) (paymentsSlice map[string]models.Payment, err error) {
}

func (service *GatewayService) performUserLoyaltyGetRequest(
	username string,
) (loyalty models.Loyalty, err error) {
}

func (service *GatewayService) performHotelByUidGetRequest(
	hotelUid string,
) (hotel models.Hotel, err error) {
}

func (service *GatewayService) performPaymentPostRequest(
	price int,
) (payment models.Payment, err error) {
}

func (service *GatewayService) performLoyaltyPatchRequest(
	loyalty *models.Loyalty,
) (err error) {
}

func (service *GatewayService) performReservationPostRequest(
	username string,
	paymentUid string,
	hotelId int,
	status string,
	startDate string,
	endDate string,
) (reservation models.Reservation, err error) {
}

func (service *GatewayService) performReservGetRequest(
	reservUid string,
) (reserv models.Reservation, err error) {
}

func (service *GatewayService) performHotelByIdGetRequest(
	hoteId int,
) (hotel models.Hotel, err error) {
}

func (service *GatewayService) performLoyaltyGetByUsernameRequest(
	username string,
) (loyalty models.Loyalty, err error) {
}

func (service *GatewayService) performReservPatchRequest(
	reservation *models.Reservation,
) (err error) {
}

func (service *GatewayService) performPaymentPatchRequest(
	payment *models.Payment,
) (err error) {
}

func (service *GatewayService) ReadAllHotels(
	page int,
	pageSize int,
) (pagRes models.PagiationResponse, err error) {
	if page <= 0 || pageSize <= 0 {
		return pagRes, errors.New(serverrors.ErrInvalidPagesData)
	}

	hotelsSlice, err := service.performAllHotelsGetRequest()

	if err != nil {
		return pagRes, err
	}

	models.HotelsSliceToPagRes(&pagRes, hotelsSlice, page, pageSize)

	return pagRes, nil
}

func (service *GatewayService) ReadUserInfo(
	username string,
) (userInfoRes models.UserInfoResponse, err error) {
	if strings.Trim(username, ` `) == `` {
		return userInfoRes, errors.New(serverrors.ErrInvalidUsername)
	}

	userReservsSlice, err := service.performUserReservsGetRequest(username)

	if err != nil {
		return userInfoRes, err
	}

	hotelsMap, err := service.performReservsHotelsGetRequest(userReservsSlice)

	if err != nil {
		return userInfoRes, err
	}

	paymentsMap, err := service.performReservsPaymentsGetRequest(userReservsSlice)

	if err != nil {
		return userInfoRes, err
	}

	loyalty, err := service.performUserLoyaltyGetRequest(username)

	if err != nil && !errors.Is(err, errors.New(serverrors.ErrEntityNotFound)) {
		return userInfoRes, err
	}

	models.ReservsSliceToUserInfoRes(&userInfoRes, userReservsSlice, hotelsMap, paymentsMap, &loyalty)

	return userInfoRes, nil
}

func (service *GatewayService) ReadUserReservations(
	username string,
) (reservsResSlice []models.ReservationResponse, err error) {
	if strings.Trim(username, ` `) == `` {
		return reservsResSlice, errors.New(serverrors.ErrInvalidUsername)
	}

	userReservsSlice, err := service.performUserReservsGetRequest(username)

	if err != nil {
		return reservsResSlice, err
	}

	hotelsMap, err := service.performReservsHotelsGetRequest(userReservsSlice)

	if err != nil {
		return reservsResSlice, err
	}

	paymentsMap, err := service.performReservsPaymentsGetRequest(userReservsSlice)

	if err != nil {
		return reservsResSlice, err
	}

	models.ReservsSliceToReservRes(reservsResSlice, userReservsSlice, hotelsMap, paymentsMap)

	return reservsResSlice, nil
}

func (service *GatewayService) CreateReservation(
	username string,
	crReservReq *models.CreateReservationRequest,
) (crReservRes models.CreateReservationResponse, err error) {
	if strings.Trim(username, ` `) == `` {
		return crReservRes, errors.New(serverrors.ErrInvalidUsername)
	}

	_, err = models.ValidateCrReservReq(crReservReq)

	if err != nil {
		return crReservRes, errors.New(serverrors.ErrInvalidCrReservReq)
	}

	hotel, err := service.performHotelByUidGetRequest(crReservReq.HotelUid)

	if err != nil {
		return crReservRes, err
	}

	loyalty, err := service.performUserLoyaltyGetRequest(username)

	if err != nil && !errors.Is(err, errors.New(serverrors.ErrEntityNotFound)) {
		return crReservRes, err
	}

	startDate, _ := time.Parse(`%F`, crReservReq.StartDate)
	endDate, _ := time.Parse(`%F`, crReservReq.EndDate)
	nightsCount := int(endDate.Sub(startDate).Hours() / 24)
	price := nightsCount * hotel.Price

	if loyalty.Discount > 0 {
		price = int(math.Round(float64(price) * float64(loyalty.Discount/100.0)))
	}

	payment, err := service.performPaymentPostRequest(price)

	if err != nil {
		return crReservRes, err
	}

	reservation, err := service.performReservationPostRequest(
		username, payment.Uid,
		hotel.Id, payment.Status,
		crReservReq.StartDate, crReservReq.EndDate,
	)

	if err != nil {
		return crReservRes, err
	}

	loyalty.ReservationCount++

	err = service.performLoyaltyPatchRequest(&loyalty)

	if err != nil {
		return crReservRes, err
	}

	models.ReservToCrReservRes(&crReservRes, &reservation, &payment, &loyalty)

	return crReservRes, nil
}

func (service *GatewayService) ReadReservation(
	reservUid string,
	username string,
) (reservRes models.ReservationResponse, err error) {
	if strings.Trim(username, ` `) == `` {
		return reservRes, errors.New(serverrors.ErrInvalidUsername)
	}

	if uuid.Validate(reservUid) != nil {
		return reservRes, errors.New(serverrors.ErrInvalidReservUid)
	}

	reservation, err := service.performReservGetRequest(reservUid)

	if err != nil {
		return reservRes, err
	}

	if reservation.Username != username {
		return reservRes, errors.New(serverrors.ErrInvalidUsername)
	}

	hotel, err := service.performHotelByIdGetRequest(reservation.HotelId)

	if err != nil {
		return reservRes, err
	}

	payment, err := service.performPaymentByUidGet(reservation.PaymentUid)

	if err != nil {
		return reservRes, err
	}

	models.ReservToReservRes(&reservRes, &reservation, &hotel, &payment)

	return reservRes, nil
}

func (service *GatewayService) DeleteReservation(
	reservUid string,
	username string,
) error {
	if strings.Trim(username, ` `) == `` {
		return errors.New(serverrors.ErrInvalidUsername)
	}

	if uuid.Validate(reservUid) != nil {
		return errors.New(serverrors.ErrInvalidReservUid)
	}

	reservation, err := service.performReservGetRequest(reservUid)

	if err != nil {
		return err
	}

	reservation.Status = `CANCELED`
	err = service.performReservPatchRequest(&reservation)

	if err != nil {
		return err
	}

	payment, err := service.performPaymentByUidGet(reservation.PaymentUid)

	if err != nil {
		return err
	}

	payment.Status = `CANCELED`
	err = service.performPaymentPatchRequest(&payment)

	if err != nil {
		return err
	}

	loyalty, err := service.performLoyaltyGetByUsernameRequest(username)

	if err != nil {
		return err
	}

	loyalty.ReservationCount--

	return service.performLoyaltyPatchRequest(&loyalty)
}

func (service *GatewayService) ReadUserLoyalty(
	username string,
) (loyaltyInfoRes models.LoyaltyInfoResponse, err error) {
	if strings.Trim(username, ` `) == `` {
		return loyaltyInfoRes, errors.New(serverrors.ErrInvalidUsername)
	}

	loyalty, err := service.performLoyaltyGetByUsernameRequest(username)

	if err != nil {
		return loyaltyInfoRes, err
	}

	models.LoyaltyToLoyaltyInfoRes(&loyaltyInfoRes, &loyalty)

	return loyaltyInfoRes, nil
}
