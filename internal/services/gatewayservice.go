package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/agarmirus/ds-lab02/internal/models"
	"github.com/google/uuid"
)

type GatewayService struct {
	host string
	port int

	reservationServiceHost string
	reservationServicePort int

	paymentServiceHost string
	paymentServicePort int

	loyaltyServiceHost string
	loyaltyServicePort int
}

func NewGatewayService(
	host string,
	port int,
	reservationServiceHost string,
	reservationServicePort int,
	paymentServiceHost string,
	paymentServicePort int,
	loyaltyServiceHost string,
	loyaltyServicePort int,
) IService {
	return &GatewayService{
		host,
		port,
		reservationServiceHost,
		reservationServicePort,
		paymentServiceHost,
		paymentServicePort,
		loyaltyServiceHost,
		loyaltyServicePort,
	}
}

func parseHotelsRes(reservServiceRes *http.Response) []models.Hotel {
	var hotelServiceResBody []byte
	reservServiceRes.Body.Read(hotelServiceResBody)

	var resultMap map[string]interface{}
	json.Unmarshal(hotelServiceResBody, &resultMap)

	totalElements := resultMap["totalElements"].(int)

	var resHotels []models.Hotel

	if totalElements > 0 {
		var itemsSlice []map[string]interface{}
		json.Unmarshal(resultMap["items"].([]byte), &itemsSlice)

		for i := 0; i < totalElements; i++ {
			hotel := models.Hotel{}
			hotel.SetId(itemsSlice[i]["id"].(int))
			hotel.SetUid(itemsSlice[i]["hotelUid"].(uuid.UUID))
			hotel.SetName(itemsSlice[i]["name"].(string))
			hotel.SetCountry(itemsSlice[i]["country"].(string))
			hotel.SetCity(itemsSlice[i]["city"].(string))
			hotel.SetAddress(itemsSlice[i]["address"].(string))
			hotel.SetStars(itemsSlice[i]["stars"].(int))
			hotel.SetPrice(itemsSlice[i]["price"].(int))

			resHotels = append(resHotels, hotel)
		}
	}

	return resHotels
}

func parseReservationsRes(res *http.Response) []models.Reservation {
	var reservationsSlice []models.Reservation

	var reservationsResBody []byte
	res.Body.Read(reservationsResBody)

	var resultMap map[string][]map[string]interface{}
	json.Unmarshal(reservationsResBody, &resultMap)

	rawReservationsSlice := resultMap["reservations"]

	for _, rawReservation := range rawReservationsSlice {
		var reservation models.Reservation
		reservation.SetId(rawReservation["id"].(int))
		reservation.SetUid(rawReservation["reservationUid"].(uuid.UUID))
		reservation.SetUsername(rawReservation["username"].(string))
		var payment models.Payment
		payment.SetUid(rawReservation["paymentUid"].(uuid.UUID))
		reservation.SetPayment(payment)
		var hotel models.Hotel
		hotel.SetId(rawReservation["hotelId"].(int))
		reservation.SetHotel(hotel)
		reservation.SetStatus(rawReservation["status"].(string))
		reservation.SetStartDate(rawReservation["startDate"].(time.Time))
		reservation.SetEndDate(rawReservation["endDate"].(time.Time))

		reservationsSlice = append(reservationsSlice, reservation)
	}

	return reservationsSlice
}

func parseLoyaltyRes(res *http.Response) models.Loyalty {
	var loyalty models.Loyalty

	var loyaltyResBody []byte
	res.Body.Read(loyaltyResBody)

	var rawLoyaltyMap map[string]interface{}
	json.Unmarshal(loyaltyResBody, &rawLoyaltyMap)

	loyalty.SetId(rawLoyaltyMap["id"].(int))
	loyalty.SetUsername(rawLoyaltyMap["username"].(string))
	loyalty.SetReservationCount(rawLoyaltyMap["reservationCount"].(int))
	loyalty.SetStatus(rawLoyaltyMap["status"].(string))
	loyalty.SetDiscount(rawLoyaltyMap["discount"].(int))

	return loyalty
}

func parseSingleHotelRes(res *http.Response) models.Hotel {
	var hotel models.Hotel

	var hotelResBody []byte
	res.Body.Read(hotelResBody)

	var rawHotelMap map[string]interface{}
	json.Unmarshal(hotelResBody, &rawHotelMap)

	hotel.SetId(rawHotelMap["id"].(int))
	hotel.SetUid(rawHotelMap["uid"].(uuid.UUID))
	hotel.SetName(rawHotelMap["name"].(string))
	hotel.SetCountry(rawHotelMap["country"].(string))
	hotel.SetCity(rawHotelMap["city"].(string))
	hotel.SetAddress(rawHotelMap["address"].(string))
	hotel.SetStars(rawHotelMap["stars"].(int))
	hotel.SetPrice(rawHotelMap["price"].(int))

	return hotel
}

func parsePaymentRes(res *http.Response) models.Payment {
	var payment models.Payment

	var paymentResBody []byte
	res.Body.Read(paymentResBody)

	var rawPaymentMap map[string]interface{}
	json.Unmarshal(paymentResBody, &rawPaymentMap)

	payment.SetId(rawPaymentMap["id"].(int))
	payment.SetUid(rawPaymentMap["uid"].(uuid.UUID))
	payment.SetStatus(rawPaymentMap["status"].(string))
	payment.SetPrice(rawPaymentMap["price"].(int))

	return payment
}

func formHotelsReadResultJson(resHotels []models.Hotel, page int, pageSize int) []byte {
	var totalElements int = len(resHotels)

	var resultMap map[string]interface{}
	resultMap["page"] = page
	resultMap["pageSize"] = pageSize
	resultMap["totalElements"] = totalElements

	if totalElements > 0 {
		var hotelsMapsSlice []map[string]interface{}

		for i := 0; i < totalElements; i++ {
			hotelPtr := &resHotels[i]

			var hotelMap map[string]interface{}
			hotelMap["hotelUid"] = hotelPtr.GetUid()
			hotelMap["name"] = hotelPtr.GetName()
			hotelMap["country"] = hotelPtr.GetCountry()
			hotelMap["city"] = hotelPtr.GetCity()
			hotelMap["address"] = hotelPtr.GetAddress()
			hotelMap["stars"] = hotelPtr.GetStars()
			hotelMap["price"] = hotelPtr.GetPrice()

			hotelsMapsSlice = append(hotelsMapsSlice, hotelMap)
		}
	}

	resultJson, _ := json.Marshal(resultMap)
	return resultJson
}

func formReservationsReadResultJSON(reservationsSlice []models.Reservation, loyalty models.Loyalty) []byte {
	totalElements := len(reservationsSlice)

	var resultMap map[string]interface{}
	var rawReservationsMapsSlice []map[string]interface{}

	for i := 0; i < totalElements; i++ {
		hotel := reservationsSlice[i].GetHotel()
		payment := reservationsSlice[i].GetPayment()

		var rawHotelMap map[string]interface{}
		rawHotelMap["hotelUid"] = hotel.GetUid()
		rawHotelMap["name"] = hotel.GetName()
		rawHotelMap["fullAddress"] = hotel.GetAddress()
		rawHotelMap["stars"] = hotel.GetStars()

		var rawPaymentMap map[string]interface{}
		rawPaymentMap["status"] = payment.GetStatus()
		rawPaymentMap["price"] = payment.GetPrice()

		var rawReservationMap map[string]interface{}
		rawReservationMap["reservationUid"] = reservationsSlice[i].GetUid()
		rawReservationMap["hotel"] = rawHotelMap
		rawReservationMap["startDate"] = reservationsSlice[i].GetStartDate()
		rawReservationMap["endDate"] = reservationsSlice[i].GetEndDate()
		rawReservationMap["status"] = reservationsSlice[i].GetStatus()
		rawReservationMap["payment"] = rawPaymentMap

		rawReservationsMapsSlice = append(rawReservationsMapsSlice, rawReservationMap)
	}

	var rawLoyaltyMap map[string]interface{}
	rawLoyaltyMap["status"] = loyalty.GetStatus()
	rawLoyaltyMap["discount"] = loyalty.GetDiscount()

	resultMap["reservations"] = rawReservationsMapsSlice
	resultMap["loyalty"] = rawLoyaltyMap

	resultJson, _ := json.Marshal(resultMap)
	return resultJson
}

func (service *GatewayService) handleHotelsRead(res http.ResponseWriter, req *http.Request) {
	page, err := strconv.Atoi(req.FormValue("page"))

	if err == nil {
		var pageSize int
		pageSize, err = strconv.Atoi(req.FormValue("size"))

		if err == nil {
			var reservServiceRes *http.Response
			reservServiceRes, err = http.Get(
				fmt.Sprintf(
					"http://%s:%d/api/v1/hotels",
					service.reservationServiceHost,
					service.reservationServicePort,
				),
			)

			if err != nil {
				res.WriteHeader(http.StatusServiceUnavailable)
			} else {
				if reservServiceRes.StatusCode == http.StatusOK {
					var resHotels []models.Hotel = parseHotelsRes(reservServiceRes)
					var resultJson []byte = formHotelsReadResultJson(resHotels, page, pageSize)

					res.WriteHeader(http.StatusOK)
					res.Write(resultJson)
				} else {
					res.WriteHeader(http.StatusServiceUnavailable)
				}
			}
		} else {
			res.WriteHeader(http.StatusBadRequest)
		}
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}

func (service *GatewayService) handleReservationsGetByUsername(username string) ([]models.Reservation, error) {
	var err error
	var reservationsSlice []models.Reservation

	reservServiceReq, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"http://%s:%d/api/v1/reservations",
			service.reservationServiceHost,
			service.reservationServicePort,
		),
		nil,
	)

	if err == nil {
		reservServiceReq.Header.Set("X-User-Name", username)

		var res *http.Response
		res, err = (&http.Client{}).Do(reservServiceReq)

		if err == nil {
			reservationsSlice = parseReservationsRes(res)
		}
	}

	return reservationsSlice, err
}

func (service *GatewayService) handleLoyaltyGetByUsername(username string) (models.Loyalty, error) {
	var err error
	var loyalty models.Loyalty

	loyaltyServiceReq, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"http://%s:%d/api/v1/loyalty",
			service.loyaltyServiceHost,
			service.loyaltyServicePort,
		),
		nil,
	)

	if err == nil {
		loyaltyServiceReq.Header.Set("X-User-Name", username)

		var res *http.Response
		res, err = (&http.Client{}).Do(loyaltyServiceReq)

		if err == nil {
			loyalty = parseLoyaltyRes(res)
		}
	}

	return loyalty, err
}

func (service *GatewayService) handleHotelGetById(id int) (models.Hotel, error) {
	var err error
	var hotel models.Hotel
	reservServiceReq, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"http://%s:%d/api/v1/hotels/%d",
			service.reservationServiceHost,
			service.reservationServiceHost,
			id,
		),
		nil,
	)

	if err == nil {
		var res *http.Response
		res, err = (&http.Client{}).Do(reservServiceReq)

		if err == nil {
			hotel = parseSingleHotelRes(res)
		}
	}

	return hotel, err
}

func (service *GatewayService) handlePaymentGetByUid(uid uuid.UUID) (models.Payment, error) {
	var err error
	var payment models.Payment
	reservServiceReq, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"http://%s:%d/api/v1/hotels/%s",
			service.reservationServiceHost,
			service.reservationServiceHost,
			uid.String(),
		),
		nil,
	)

	if err == nil {
		var res *http.Response
		res, err = (&http.Client{}).Do(reservServiceReq)

		if err == nil {
			payment = parsePaymentRes(res)
		}
	}

	return payment, err
}

func (service *GatewayService) handleHotelsGetByReservations(reservationsSlice []models.Reservation) ([]models.Reservation, error) {
	var err error = nil
	reservationsCount := len(reservationsSlice)

	for i := 0; i < reservationsCount; i++ {
		var hotel models.Hotel
		hotel, err = service.handleHotelGetById(reservationsSlice[i].GetHotel().GetId())

		if err == nil {
			reservationsSlice[i].SetHotel(hotel)
		} else {
			i = reservationsCount
		}
	}

	return reservationsSlice, err
}

func (service *GatewayService) handlePaymentsGetByReservations(reservationsSlice []models.Reservation) ([]models.Reservation, error) {
	var err error = nil
	reservationsCount := len(reservationsSlice)

	for i := 0; i < reservationsCount; i++ {
		var payment models.Payment
		payment, err = service.handlePaymentGetByUid(reservationsSlice[i].GetPayment().GetUid())

		if err == nil {
			reservationsSlice[i].SetPayment(payment)
		} else {
			i = reservationsCount
		}
	}

	return reservationsSlice, err
}

func (service *GatewayService) handleUserRead(res http.ResponseWriter, req *http.Request) {
	var err error

	username := strings.Trim(req.Header.Get("X-User-Name"), " ")

	reservationsSlice, err := service.handleReservationsGetByUsername(username)

	if err == nil && len(reservationsSlice) > 0 {
		var loyalty models.Loyalty
		loyalty, err = service.handleLoyaltyGetByUsername(username)

		if err == nil {
			reservationsSlice, err = service.handleHotelsGetByReservations(reservationsSlice)

			if err == nil {
				reservationsSlice, err = service.handlePaymentsGetByReservations(reservationsSlice)

				if err == nil {
					resultJSON := formReservationsReadResultJSON(reservationsSlice, loyalty)

					res.WriteHeader(http.StatusOK)
					res.Write(resultJSON)
				} else {
					res.WriteHeader(http.StatusServiceUnavailable)
				}
			} else {
				res.WriteHeader(http.StatusServiceUnavailable)
			}
		} else {
			res.WriteHeader(http.StatusServiceUnavailable)
		}
	} else if err == nil {
		res.WriteHeader(http.StatusNotFound)
	} else {
		res.WriteHeader(http.StatusServiceUnavailable)
	}
}

func (service *GatewayService) handleHotelsRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		service.handleHotelsRead(res, req)
	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (service *GatewayService) handleUserRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		service.handleUserRead(res, req)
	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (service *GatewayService) handleReservationsRequest(res http.ResponseWriter, req *http.Request) {
	// TODO
}

func (service *GatewayService) handleReservationUidRequest(res http.ResponseWriter, req *http.Request) {
	// TODO
}

func (service *GatewayService) handleLoyaltyRequest(res http.ResponseWriter, req *http.Request) {
	// TODO
}

func (service *GatewayService) Prepare() error {
	http.HandleFunc("/api/v1/hotels", service.handleHotelsRequest)
	http.HandleFunc("/api/v1/me", service.handleUserRequest)
	http.HandleFunc("/api/v1/reservations/{reservationUid}", service.handleReservationUidRequest)
	http.HandleFunc("/api/v1/reservations", service.handleReservationsRequest)
	http.HandleFunc("/api/v1/loyalty", service.handleLoyaltyRequest)

	return nil
}

func (service *GatewayService) Run() error {
	return http.ListenAndServe(fmt.Sprintf("%s:%d", service.host, service.port), nil)
}
