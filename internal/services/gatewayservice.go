package services

// TODO: проверить коды ответов

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
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
	loyaltyServicePort int
}

func NewGatewayService(
	reservServiceHost string,
	reservServicePort int,
	paymentServiceHost string,
	paymentServicePort int,
	loyaltyServiceHost string,
	loyaltyServicePort int,
) IGatewayService {
	return &GatewayService{
		reservServiceHost,
		reservServicePort,
		paymentServiceHost,
		paymentServicePort,
		loyaltyServiceHost,
		loyaltyServicePort,
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
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"http://%s:%d/api/v1/reservations",
			service.reservServiceHost,
			service.reservServicePort,
		),
		nil,
	)

	if err != nil {
		return userReservsSlice, errors.New(serverrors.ErrNewRequestForming)
	}

	req.Header.Add(`X-User-Name`, username)
	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return userReservsSlice, errors.New(serverrors.ErrRequestSend)
	}

	var resBody []byte
	_, err = res.Body.Read(resBody)

	if err != nil {
		return userReservsSlice, errors.New(serverrors.ErrResponseRead)
	}

	err = json.Unmarshal(resBody, &userReservsSlice)

	if err != nil {
		return userReservsSlice, errors.New(serverrors.ErrResponseParse)
	}

	return userReservsSlice, nil
}

func (service *GatewayService) performHotelByIdGetRequest(
	hotelId int,
) (hotel models.Hotel, err error) {
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
		return hotel, errors.New(serverrors.ErrNewRequestForming)
	}

	req.Header.Add(`Hotel-Id`, strconv.Itoa(hotelId))
	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return hotel, errors.New(serverrors.ErrRequestSend)
	}

	var resBody []byte
	_, err = res.Body.Read(resBody)

	if err != nil {
		return hotel, errors.New(serverrors.ErrResponseRead)
	}

	err = json.Unmarshal(resBody, &hotel)

	if err != nil {
		return hotel, errors.New(serverrors.ErrResponseParse)
	}

	return hotel, nil
}

func (service *GatewayService) performReservsHotelsGetRequest(
	userReservsSlice []models.Reservation,
) (hotelsMap map[int]models.Hotel, err error) {
	hotelsMap = make(map[int]models.Hotel)
	totalElements := len(userReservsSlice)

	for i := 0; i < totalElements; i++ {
		hotel, err := service.performHotelByIdGetRequest(userReservsSlice[i].HotelId)

		if err != nil {
			return hotelsMap, err
		}

		hotelsMap[hotel.Id] = hotel
	}

	return hotelsMap, nil
}

func (service *GatewayService) performPaymentByUidGetRequest(
	paymentUid string,
) (payment models.Payment, err error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"http://%s:%d/api/v1/payment/%s",
			service.paymentServiceHost,
			service.paymentServicePort,
			paymentUid,
		),
		nil,
	)

	if err != nil {
		return payment, errors.New(serverrors.ErrNewRequestForming)
	}

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return payment, errors.New(serverrors.ErrRequestSend)
	}

	var resBody []byte
	_, err = res.Body.Read(resBody)

	if err != nil {
		return payment, errors.New(serverrors.ErrResponseRead)
	}

	err = json.Unmarshal(resBody, &payment)

	if err != nil {
		return payment, errors.New(serverrors.ErrResponseParse)
	}

	return payment, nil
}

func (service *GatewayService) performReservsPaymentsGetRequest(
	userReservsSlice []models.Reservation,
) (paymentsMap map[string]models.Payment, err error) {
	paymentsMap = make(map[string]models.Payment)
	totalElements := len(userReservsSlice)

	for i := 0; i < totalElements; i++ {
		payment, err := service.performPaymentByUidGetRequest(userReservsSlice[i].PaymentUid)

		if err != nil {
			return paymentsMap, err
		}

		paymentsMap[payment.Uid] = payment
	}

	return paymentsMap, nil
}

func (service *GatewayService) performUserLoyaltyGetRequest(
	username string,
) (loyalty models.Loyalty, err error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"http://%s:%d/api/v1/loyalty",
			service.loyaltyServiceHost,
			service.loyaltyServicePort,
		),
		nil,
	)

	if err != nil {
		return loyalty, errors.New(serverrors.ErrNewRequestForming)
	}

	req.Header.Add(`X-User-Name`, username)
	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return loyalty, errors.New(serverrors.ErrRequestSend)
	}

	var resBody []byte
	_, err = res.Body.Read(resBody)

	if err != nil {
		return loyalty, errors.New(serverrors.ErrResponseRead)
	}

	err = json.Unmarshal(resBody, &loyalty)

	if err != nil {
		return loyalty, errors.New(serverrors.ErrResponseParse)
	}

	return loyalty, nil
}

func (service *GatewayService) performHotelByUidGetRequest(
	hotelUid string,
) (hotel models.Hotel, err error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"http://%s:%d/api/v1/hotels/%s",
			service.reservServiceHost,
			service.reservServicePort,
			hotelUid,
		),
		nil,
	)

	if err != nil {
		return hotel, errors.New(serverrors.ErrNewRequestForming)
	}

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return hotel, errors.New(serverrors.ErrRequestSend)
	}

	var resBody []byte
	_, err = res.Body.Read(resBody)

	if err != nil {
		return hotel, errors.New(serverrors.ErrResponseRead)
	}

	err = json.Unmarshal(resBody, &hotel)

	if err != nil {
		return hotel, errors.New(serverrors.ErrResponseParse)
	}

	return hotel, nil
}

func (service *GatewayService) performPaymentPostRequest(
	price int,
) (payment models.Payment, err error) {
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(
			"http://%s:%d/api/v1/payment",
			service.reservServiceHost,
			service.reservServicePort,
		),
		nil,
	)

	if err != nil {
		return payment, errors.New(serverrors.ErrNewRequestForming)
	}

	req.Header.Add(`Price`, strconv.Itoa(price))
	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return payment, errors.New(serverrors.ErrRequestSend)
	}

	var resBody []byte
	_, err = res.Body.Read(resBody)

	if err != nil {
		return payment, errors.New(serverrors.ErrResponseRead)
	}

	err = json.Unmarshal(resBody, &payment)

	if err != nil {
		return payment, errors.New(serverrors.ErrResponseParse)
	}

	return payment, nil
}

func (service *GatewayService) performLoyaltyPatchRequest(
	loyalty *models.Loyalty,
) (err error) {
	loyaltyJSON, err := json.Marshal(loyalty)

	if err != nil {
		return errors.New(serverrors.ErrJSONParse)
	}

	req, err := http.NewRequest(
		"PATCH",
		fmt.Sprintf(
			"http://%s:%d/api/v1/loyalty/%d",
			service.loyaltyServiceHost,
			service.loyaltyServicePort,
			loyalty.Id,
		),
		bytes.NewBuffer(loyaltyJSON),
	)

	if err != nil {
		return errors.New(serverrors.ErrNewRequestForming)
	}

	_, err = (&http.Client{}).Do(req)

	if err != nil {
		return errors.New(serverrors.ErrRequestSend)
	}

	return nil
}

func (service *GatewayService) performReservationPostRequest(
	username string,
	paymentUid string,
	hotelId int,
	status string,
	startDate string,
	endDate string,
) (reservation models.Reservation, err error) {
	var newReservation models.Reservation
	newReservation.Username = username
	newReservation.PaymentUid = paymentUid
	newReservation.HotelId = hotelId
	newReservation.Status = status
	newReservation.StartDate, _ = time.Parse(`%F`, startDate)
	newReservation.EndDate, _ = time.Parse(`%F`, endDate)

	newReservJSON, err := json.Marshal(newReservation)

	if err != nil {
		return reservation, errors.New(serverrors.ErrJSONParse)
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(
			"http://%s:%d/api/v1/payment",
			service.reservServiceHost,
			service.reservServicePort,
		),
		bytes.NewBuffer(newReservJSON),
	)

	if err != nil {
		return reservation, errors.New(serverrors.ErrNewRequestForming)
	}

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return reservation, errors.New(serverrors.ErrRequestSend)
	}

	var resBody []byte
	_, err = res.Body.Read(resBody)

	if err != nil {
		return reservation, errors.New(serverrors.ErrResponseRead)
	}

	err = json.Unmarshal(resBody, &reservation)

	if err != nil {
		return reservation, errors.New(serverrors.ErrResponseParse)
	}

	return reservation, nil
}

func (service *GatewayService) performReservGetRequest(
	reservUid string,
) (reserv models.Reservation, err error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"http://%s:%d/api/v1/reservations/%s",
			service.reservServiceHost,
			service.reservServicePort,
			reservUid,
		),
		nil,
	)

	if err != nil {
		return reserv, errors.New(serverrors.ErrNewRequestForming)
	}

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return reserv, errors.New(serverrors.ErrRequestSend)
	}

	var resBody []byte
	_, err = res.Body.Read(resBody)

	if err != nil {
		return reserv, errors.New(serverrors.ErrResponseRead)
	}

	err = json.Unmarshal(resBody, &reserv)

	if err != nil {
		return reserv, errors.New(serverrors.ErrResponseParse)
	}

	return reserv, nil
}

func (service *GatewayService) performLoyaltyGetByUsernameRequest(
	username string,
) (loyalty models.Loyalty, err error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"http://%s:%d/api/v1/loyalty",
			service.loyaltyServiceHost,
			service.loyaltyServicePort,
		),
		nil,
	)

	if err != nil {
		return loyalty, errors.New(serverrors.ErrNewRequestForming)
	}

	req.Header.Add(`X-User-Name`, username)
	res, err := (&http.Client{}).Do(req)

	if err != nil {
		return loyalty, errors.New(serverrors.ErrRequestSend)
	}

	var resBody []byte
	_, err = res.Body.Read(resBody)

	if err != nil {
		return loyalty, errors.New(serverrors.ErrResponseRead)
	}

	err = json.Unmarshal(resBody, &loyalty)

	if err != nil {
		return loyalty, errors.New(serverrors.ErrResponseParse)
	}

	return loyalty, nil
}

func (service *GatewayService) performReservPatchRequest(
	reservation *models.Reservation,
) (err error) {
	reservJSON, err := json.Marshal(reservation)

	if err != nil {
		return errors.New(serverrors.ErrJSONParse)
	}

	req, err := http.NewRequest(
		"PATCH",
		fmt.Sprintf(
			"http://%s:%d/api/v1/reservations/%s",
			service.loyaltyServiceHost,
			service.loyaltyServicePort,
			reservation.Uid,
		),
		bytes.NewBuffer(reservJSON),
	)

	if err != nil {
		return errors.New(serverrors.ErrNewRequestForming)
	}

	_, err = (&http.Client{}).Do(req)

	if err != nil {
		return errors.New(serverrors.ErrRequestSend)
	}

	return nil
}

func (service *GatewayService) performPaymentPatchRequest(
	payment *models.Payment,
) (err error) {
	paymentJSON, err := json.Marshal(payment)

	if err != nil {
		return errors.New(serverrors.ErrJSONParse)
	}

	req, err := http.NewRequest(
		"PATCH",
		fmt.Sprintf(
			"http://%s:%d/api/v1/payment/%s",
			service.paymentServiceHost,
			service.paymentServicePort,
			payment.Uid,
		),
		bytes.NewBuffer(paymentJSON),
	)

	if err != nil {
		return errors.New(serverrors.ErrNewRequestForming)
	}

	_, err = (&http.Client{}).Do(req)

	if err != nil {
		return errors.New(serverrors.ErrRequestSend)
	}

	return nil
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

	payment, err := service.performPaymentByUidGetRequest(reservation.PaymentUid)

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

	payment, err := service.performPaymentByUidGetRequest(reservation.PaymentUid)

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
