package services

// TODO: проверить коды ответов

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

func (service *GatewayService) performAllHotelsGetRequest(
	page int,
	pageSize int,
) (hotelsSlice []models.Hotel, err error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"http://%s:%d/api/v1/hotels?page=%d&size=%d",
			service.reservServiceHost,
			service.reservServicePort,
			page,
			pageSize,
		),
		nil,
	)

	if err != nil {
		log.Println("[ERROR] GatewayService.performAllHotelsGetRequest. Error while creating new HTTP-request: ", err)
		return hotelsSlice, errors.New(serverrors.ErrNewRequestForming)
	}

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		log.Println("[ERROR] GatewayService.performAllHotelsGetRequest. Error while sending request")
		return hotelsSlice, errors.New(serverrors.ErrRequestSend)
	}

	var resBody []byte
	_, err = res.Body.Read(resBody)

	if err != nil {
		log.Println("[ERROR] GatewayService.performAllHotelsGetRequest. Error while reading response")
		return hotelsSlice, errors.New(serverrors.ErrResponseRead)
	}

	err = json.Unmarshal(resBody, &hotelsSlice)

	if err != nil {
		log.Println("[ERROR] GatewayService.performAllHotelsGetRequest. Error while parsing JSON response body")
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
		log.Println("[ERROR] GatewayService.performUserReservsGetRequest. Error while creating new request")
		return userReservsSlice, errors.New(serverrors.ErrNewRequestForming)
	}

	req.Header.Add(`X-User-Name`, username)
	res, err := (&http.Client{}).Do(req)

	if err != nil {
		log.Println("[ERROR] GatewayService.performUserReservsGetRequest. Error while sending request")
		return userReservsSlice, errors.New(serverrors.ErrRequestSend)
	}

	var resBody []byte
	_, err = res.Body.Read(resBody)

	if err != nil {
		log.Println("[ERROR] GatewayService.performUserReservsGetRequest. Error while reading response")
		return userReservsSlice, errors.New(serverrors.ErrResponseRead)
	}

	err = json.Unmarshal(resBody, &userReservsSlice)

	if err != nil {
		log.Println("[ERROR] GatewayService.performUserReservsGetRequest. Error while parsing JSON response body")
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
		log.Println("[ERROR] GatewayService.performHotelByIdGetRequest. Error while creating new request")
		return hotel, errors.New(serverrors.ErrNewRequestForming)
	}

	req.Header.Add(`Hotel-Id`, strconv.Itoa(hotelId))
	res, err := (&http.Client{}).Do(req)

	if err != nil {
		log.Println("[ERROR] GatewayService.performHotelByIdGetRequest. Error while sending request")
		return hotel, errors.New(serverrors.ErrRequestSend)
	}

	var resBody []byte
	_, err = res.Body.Read(resBody)

	if err != nil {
		log.Println("[ERROR] GatewayService.performHotelByIdGetRequest. Error while reading response")
		return hotel, errors.New(serverrors.ErrResponseRead)
	}

	err = json.Unmarshal(resBody, &hotel)

	if err != nil {
		log.Println("[ERROR] GatewayService.performHotelByIdGetRequest. Error while parsing JSON response body")
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
			log.Println("[ERROR] GatewayService.performReservsHotelsGetRequest. performHotelByIdGetRequest returned error:", err)
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
		log.Println("[ERROR] GatewayService.performPaymentByUidGetRequest. Error while creating new request")
		return payment, errors.New(serverrors.ErrNewRequestForming)
	}

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		log.Println("[ERROR] GatewayService.performPaymentByUidGetRequest. Error while sending request")
		return payment, errors.New(serverrors.ErrRequestSend)
	}

	var resBody []byte
	_, err = res.Body.Read(resBody)

	if err != nil {
		log.Println("[ERROR] GatewayService.performPaymentByUidGetRequest. Error while reading response")
		return payment, errors.New(serverrors.ErrResponseRead)
	}

	err = json.Unmarshal(resBody, &payment)

	if err != nil {
		log.Println("[ERROR] GatewayService.performPaymentByUidGetRequest. Error while parsing JSON response body")
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
			log.Println("[ERROR] GatewayService.performReservsPaymentsGetRequest. performPaymentByUidGetRequest returned error:", err)
			return paymentsMap, err
		}

		paymentsMap[payment.Uid] = payment
	}

	return paymentsMap, nil
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
		log.Println("[ERROR] GatewayService.performHotelByUidGetRequest. Error while creating new request")
		return hotel, errors.New(serverrors.ErrNewRequestForming)
	}

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		log.Println("[ERROR] GatewayService.performHotelByUidGetRequest. Error while sending request")
		return hotel, errors.New(serverrors.ErrRequestSend)
	}

	var resBody []byte
	_, err = res.Body.Read(resBody)

	if err != nil {
		log.Println("[ERROR] GatewayService.performHotelByUidGetRequest. Error while reading response")
		return hotel, errors.New(serverrors.ErrResponseRead)
	}

	err = json.Unmarshal(resBody, &hotel)

	if err != nil {
		log.Println("[ERROR] GatewayService.performHotelByUidGetRequest. Error while parsing JSON response body")
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
		log.Println("[ERROR] GatewayService.performPaymentPostRequest. Error while creating new request")
		return payment, errors.New(serverrors.ErrNewRequestForming)
	}

	req.Header.Add(`Price`, strconv.Itoa(price))
	res, err := (&http.Client{}).Do(req)

	if err != nil {
		log.Println("[ERROR] GatewayService.performPaymentPostRequest. Error while sending request")
		return payment, errors.New(serverrors.ErrRequestSend)
	}

	var resBody []byte
	_, err = res.Body.Read(resBody)

	if err != nil {
		log.Println("[ERROR] GatewayService.performPaymentPostRequest. Error while reading response")
		return payment, errors.New(serverrors.ErrResponseRead)
	}

	err = json.Unmarshal(resBody, &payment)

	if err != nil {
		log.Println("[ERROR] GatewayService.performPaymentPostRequest. Error while parsing JSON response body")
		return payment, errors.New(serverrors.ErrResponseParse)
	}

	return payment, nil
}

func (service *GatewayService) performLoyaltyPutRequest(
	loyalty *models.Loyalty,
) (err error) {
	loyaltyJSON, err := json.Marshal(loyalty)

	if err != nil {
		log.Println("[ERROR] GatewayService.performLoyaltyPutRequest. Cannot create JSON object for request body")
		return errors.New(serverrors.ErrJSONParse)
	}

	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf(
			"http://%s:%d/api/v1/loyalty/%d",
			service.loyaltyServiceHost,
			service.loyaltyServicePort,
			loyalty.Id,
		),
		bytes.NewBuffer(loyaltyJSON),
	)

	if err != nil {
		log.Println("[ERROR] GatewayService.performLoyaltyPutRequest. Error while creating new request")
		return errors.New(serverrors.ErrNewRequestForming)
	}

	_, err = (&http.Client{}).Do(req)

	if err != nil {
		log.Println("[ERROR] GatewayService.performLoyaltyPutRequest. Error while sending request")
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
		log.Println("[ERROR] GatewayService.performReservationPostRequest. Cannot create JSON object for request body")
		return reservation, errors.New(serverrors.ErrJSONParse)
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(
			"http://%s:%d/api/v1/reservations",
			service.reservServiceHost,
			service.reservServicePort,
		),
		bytes.NewBuffer(newReservJSON),
	)

	if err != nil {
		log.Println("[ERROR] GatewayService.performReservationPostRequest. Error while creating new request")
		return reservation, errors.New(serverrors.ErrNewRequestForming)
	}

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		log.Println("[ERROR] GatewayService.performReservationPostRequest. Error while sending request")
		return reservation, errors.New(serverrors.ErrRequestSend)
	}

	var resBody []byte
	_, err = res.Body.Read(resBody)

	if err != nil {
		log.Println("[ERROR] GatewayService.performReservationPostRequest. Error while reading response")
		return reservation, errors.New(serverrors.ErrResponseRead)
	}

	err = json.Unmarshal(resBody, &reservation)

	if err != nil {
		log.Println("[ERROR] GatewayService.performReservationPostRequest. Error while parsing JSON response body")
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
		log.Println("[ERROR] GatewayService.performReservGetRequest. Error while creating new request")
		return reserv, errors.New(serverrors.ErrNewRequestForming)
	}

	res, err := (&http.Client{}).Do(req)

	if err != nil {
		log.Println("[ERROR] GatewayService.performReservGetRequest. Error while sending request")
		return reserv, errors.New(serverrors.ErrRequestSend)
	}

	var resBody []byte
	_, err = res.Body.Read(resBody)

	if err != nil {
		log.Println("[ERROR] GatewayService.performReservGetRequest. Error while reading response")
		return reserv, errors.New(serverrors.ErrResponseRead)
	}

	err = json.Unmarshal(resBody, &reserv)

	if err != nil {
		log.Println("[ERROR] GatewayService.performReservGetRequest. Error while parsing JSON response body")
		return reserv, errors.New(serverrors.ErrResponseParse)
	}

	return reserv, nil
}

func (service *GatewayService) performLoyaltyByUsernameGetRequest(
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
		log.Println("[ERROR] GatewayService.performLoyaltyByUsernameGetRequest. Error while creating new request")
		return loyalty, errors.New(serverrors.ErrNewRequestForming)
	}

	req.Header.Add(`X-User-Name`, username)
	res, err := (&http.Client{}).Do(req)

	if err != nil {
		log.Println("[ERROR] GatewayService.performLoyaltyByUsernameGetRequest. Error while sending request")
		return loyalty, errors.New(serverrors.ErrRequestSend)
	}

	var resBody []byte
	_, err = res.Body.Read(resBody)

	if err != nil {
		log.Println("[ERROR] GatewayService.performLoyaltyByUsernameGetRequest. Error while reading response")
		return loyalty, errors.New(serverrors.ErrResponseRead)
	}

	err = json.Unmarshal(resBody, &loyalty)

	if err != nil {
		log.Println("[ERROR] GatewayService.performLoyaltyByUsernameGetRequest. Error while parsing JSON response body")
		return loyalty, errors.New(serverrors.ErrResponseParse)
	}

	return loyalty, nil
}

func (service *GatewayService) performReservPutRequest(
	reservation *models.Reservation,
) (err error) {
	reservJSON, err := json.Marshal(reservation)

	if err != nil {
		log.Println("[ERROR] GatewayService.performReservPutRequest. Cannot create JSON object for request body")
		return errors.New(serverrors.ErrJSONParse)
	}

	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf(
			"http://%s:%d/api/v1/reservations/%s",
			service.loyaltyServiceHost,
			service.loyaltyServicePort,
			reservation.Uid,
		),
		bytes.NewBuffer(reservJSON),
	)

	if err != nil {
		log.Println("[ERROR] GatewayService.performReservPutRequest. Error while creating new request")
		return errors.New(serverrors.ErrNewRequestForming)
	}

	_, err = (&http.Client{}).Do(req)

	if err != nil {
		log.Println("[ERROR] GatewayService.performReservPutRequest. Error while sending request")
		return errors.New(serverrors.ErrRequestSend)
	}

	return nil
}

func (service *GatewayService) performPaymentPutRequest(
	payment *models.Payment,
) (err error) {
	paymentJSON, err := json.Marshal(payment)

	if err != nil {
		log.Println("[ERROR] GatewayService.performPaymentPutRequest. Cannot create JSON object for request body")
		return errors.New(serverrors.ErrJSONParse)
	}

	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf(
			"http://%s:%d/api/v1/payment/%s",
			service.paymentServiceHost,
			service.paymentServicePort,
			payment.Uid,
		),
		bytes.NewBuffer(paymentJSON),
	)

	if err != nil {
		log.Println("[ERROR] GatewayService.performPaymentPutRequest. Error while creating new request")
		return errors.New(serverrors.ErrNewRequestForming)
	}

	_, err = (&http.Client{}).Do(req)

	if err != nil {
		log.Println("[ERROR] GatewayService.performPaymentPutRequest. Error while sending request")
		return errors.New(serverrors.ErrRequestSend)
	}

	return nil
}

func (service *GatewayService) ReadAllHotels(
	page int,
	pageSize int,
) (pagRes models.PagiationResponse, err error) {
	if page <= 0 || pageSize <= 0 {
		log.Println("[ERROR] GatewayService.ReadAllHotels. Invalid parameters")
		return pagRes, errors.New(serverrors.ErrInvalidPagesData)
	}

	hotelsSlice, err := service.performAllHotelsGetRequest(page, pageSize)

	if err != nil {
		log.Println("[ERROR] GatewayService.ReadAllHotels. performAllHotelsGetRequest returned error:", err)
		return pagRes, err
	}

	models.HotelsSliceToPagRes(&pagRes, hotelsSlice, page, pageSize)

	return pagRes, nil
}

func (service *GatewayService) ReadUserInfo(
	username string,
) (userInfoRes models.UserInfoResponse, err error) {
	if strings.Trim(username, ` `) == `` {
		log.Println("[ERROR] GatewayService.ReadUserInfo. Invalid username")
		return userInfoRes, errors.New(serverrors.ErrInvalidUsername)
	}

	userReservsSlice, err := service.performUserReservsGetRequest(username)

	if err != nil {
		log.Println("[ERROR] GatewayService.ReadUserInfo. performUserReservsGetRequest returned error:", err)
		return userInfoRes, err
	}

	hotelsMap, err := service.performReservsHotelsGetRequest(userReservsSlice)

	if err != nil {
		log.Println("[ERROR] GatewayService.ReadUserInfo. performReservsHotelsGetRequest returned error:", err)
		return userInfoRes, err
	}

	paymentsMap, err := service.performReservsPaymentsGetRequest(userReservsSlice)

	if err != nil {
		log.Println("[ERROR] GatewayService.ReadUserInfo. performReservsPaymentsGetRequest returned error:", err)
		return userInfoRes, err
	}

	loyalty, err := service.performLoyaltyByUsernameGetRequest(username)

	if err != nil && !errors.Is(err, errors.New(serverrors.ErrEntityNotFound)) {
		log.Println("[ERROR] GatewayService.ReadUserInfo. performLoyaltyByUsernameGetRequest returned error:", err)
		return userInfoRes, err
	}

	models.ReservsSliceToUserInfoRes(&userInfoRes, userReservsSlice, hotelsMap, paymentsMap, &loyalty)

	return userInfoRes, nil
}

func (service *GatewayService) ReadUserReservations(
	username string,
) (reservsResSlice []models.ReservationResponse, err error) {
	if strings.Trim(username, ` `) == `` {
		log.Println("[ERROR] GatewayService.ReadUserReservations. Invalid username")
		return reservsResSlice, errors.New(serverrors.ErrInvalidUsername)
	}

	userReservsSlice, err := service.performUserReservsGetRequest(username)

	if err != nil {
		log.Println("[ERROR] GatewayService.ReadUserReservations. performUserReservsGetRequest returned error:", err)
		return reservsResSlice, err
	}

	hotelsMap, err := service.performReservsHotelsGetRequest(userReservsSlice)

	if err != nil {
		log.Println("[ERROR] GatewayService.ReadUserReservations. performReservsHotelsGetRequest returned error:", err)
		return reservsResSlice, err
	}

	paymentsMap, err := service.performReservsPaymentsGetRequest(userReservsSlice)

	if err != nil {
		log.Println("[ERROR] GatewayService.ReadUserReservations. performReservsPaymentsGetRequest returned error:", err)
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
		log.Println("[ERROR] GatewayService.CreateReservation. Invalid username")
		return crReservRes, errors.New(serverrors.ErrInvalidUsername)
	}

	_, err = models.ValidateCrReservReq(crReservReq)

	if err != nil {
		log.Println("[ERROR] GatewayService.CreateReservation. Invalid create reservation request")
		return crReservRes, errors.New(serverrors.ErrInvalidCrReservReq)
	}

	hotel, err := service.performHotelByUidGetRequest(crReservReq.HotelUid)

	if err != nil {
		log.Println("[ERROR] GatewayService.CreateReservation. performHotelByUidGetRequest returned error:", err)
		return crReservRes, err
	}

	loyalty, err := service.performLoyaltyByUsernameGetRequest(username)

	if err != nil && !errors.Is(err, errors.New(serverrors.ErrEntityNotFound)) {
		log.Println("[ERROR] GatewayService.CreateReservation. performLoyaltyByUsernameGetRequest returned error:", err)
		return crReservRes, err
	}

	startDate, _ := time.Parse(`%F`, crReservReq.StartDate)
	endDate, _ := time.Parse(`%F`, crReservReq.EndDate)
	nightsCount := int(endDate.Sub(startDate).Hours() / 24)
	price := nightsCount * hotel.Price

	if loyalty.Discount > 0 {
		price = int(math.Round(float64(price) * float64(1.0-loyalty.Discount/100.0)))
	}

	payment, err := service.performPaymentPostRequest(price)

	if err != nil {
		log.Println("[ERROR] GatewayService.CreateReservation. performPaymentPostRequest returned error:", err)
		return crReservRes, err
	}

	reservation, err := service.performReservationPostRequest(
		username, payment.Uid,
		hotel.Id, payment.Status,
		crReservReq.StartDate, crReservReq.EndDate,
	)

	if err != nil {
		log.Println("[ERROR] GatewayService.CreateReservation. performReservationPostRequest returned error:", err)
		return crReservRes, err
	}

	loyalty.ReservationCount++

	err = service.performLoyaltyPutRequest(&loyalty)

	if err != nil {
		log.Println("[ERROR] GatewayService.CreateReservation. performLoyaltyPutRequest returned error:", err)
		return crReservRes, err
	}

	models.ReservToCrReservRes(&crReservRes, &reservation, &payment, &loyalty, hotel.Uid)

	return crReservRes, nil
}

func (service *GatewayService) ReadReservation(
	reservUid string,
	username string,
) (reservRes models.ReservationResponse, err error) {
	if strings.Trim(username, ` `) == `` {
		log.Println("[ERROR] GatewayService.ReadReservation. Invalid username")
		return reservRes, errors.New(serverrors.ErrInvalidUsername)
	}

	if uuid.Validate(reservUid) != nil {
		log.Println("[ERROR] GatewayService.ReadReservation. Invalid reservation uid")
		return reservRes, errors.New(serverrors.ErrInvalidReservUid)
	}

	reservation, err := service.performReservGetRequest(reservUid)

	if err != nil {
		log.Println("[ERROR] GatewayService.ReadReservation. performReservGetRequest returned error:", err)
		return reservRes, err
	}

	if reservation.Username != username {
		log.Println("[ERROR] GatewayService.ReadReservation. Invalid username")
		return reservRes, errors.New(serverrors.ErrInvalidUsername)
	}

	hotel, err := service.performHotelByIdGetRequest(reservation.HotelId)

	if err != nil {
		log.Println("[ERROR] GatewayService.ReadReservation. performHotelByIdGetRequest returned error:", err)
		return reservRes, err
	}

	payment, err := service.performPaymentByUidGetRequest(reservation.PaymentUid)

	if err != nil {
		log.Println("[ERROR] GatewayService.ReadReservation. performPaymentByUidGetRequest returned error:", err)
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
		log.Println("[ERROR] GatewayService.DeleteReservation. Invalid username")
		return errors.New(serverrors.ErrInvalidUsername)
	}

	if uuid.Validate(reservUid) != nil {
		log.Println("[ERROR] GatewayService.DeleteReservation. Invalid reservation uid")
		return errors.New(serverrors.ErrInvalidReservUid)
	}

	reservation, err := service.performReservGetRequest(reservUid)

	if err != nil {
		log.Println("[ERROR] GatewayService.DeleteReservation. Error while getting reservation by uid: ", err)
		return err
	}

	reservation.Status = `CANCELED`
	err = service.performReservPutRequest(&reservation)

	if err != nil {
		log.Println("[ERROR] GatewayService.DeleteReservation. Error while puting reservation: ", err)
		return err
	}

	payment, err := service.performPaymentByUidGetRequest(reservation.PaymentUid)

	if err != nil {
		log.Println("[ERROR] GatewayService.DeleteReservation. Error while getting payment by uid: ", err)
		return err
	}

	payment.Status = `CANCELED`
	err = service.performPaymentPutRequest(&payment)

	if err != nil {
		log.Println("[ERROR] GatewayService.DeleteReservation. Error while puting payment: ", err)
		return err
	}

	loyalty, err := service.performLoyaltyByUsernameGetRequest(username)

	if err != nil {
		log.Println("[ERROR] GatewayService.DeleteReservation. Error while getting loyalty by username: ", err)
		return err
	}

	loyalty.ReservationCount--

	err = service.performLoyaltyPutRequest(&loyalty)

	if err != nil {
		log.Println("[ERROR] GatewayService.DeleteReservation. Error while puting loyalty: ", err)
	}

	return err
}

func (service *GatewayService) ReadUserLoyalty(
	username string,
) (loyaltyInfoRes models.LoyaltyInfoResponse, err error) {
	if strings.Trim(username, ` `) == `` {
		log.Println("[ERROR] GatewayService.ReadUserLoyalty. Invalid username")
		return loyaltyInfoRes, errors.New(serverrors.ErrInvalidUsername)
	}

	loyalty, err := service.performLoyaltyByUsernameGetRequest(username)

	if err != nil {
		log.Println("[ERROR] GatewayService.ReadUserLoyalty. error while getting loyalty by username: ", err)
		return loyaltyInfoRes, err
	}

	models.LoyaltyToLoyaltyInfoRes(&loyaltyInfoRes, &loyalty)

	return loyaltyInfoRes, nil
}
