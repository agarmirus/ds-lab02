package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

func parseHotelsRes(hotelServiceRes *http.Response) []models.Hotel {
	var hotelServiceResBody []byte
	hotelServiceRes.Body.Read(hotelServiceResBody)

	var resultMap map[string]interface{}
	json.Unmarshal(hotelServiceResBody, &resultMap)

	totalElements := resultMap["totalElements"].(int)

	var resHotels []models.Hotel

	if totalElements > 0 {
		var itemsSlice []map[string]interface{}
		json.Unmarshal(resultMap["items"].([]byte), &itemsSlice)

		for i := 0; i < totalElements; i++ {
			hotel := models.Hotel{}
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

func formHotelsResultJson(resHotels []models.Hotel, page int, pageSize int) []byte {
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

func (service *GatewayService) handleHotelsRead(res http.ResponseWriter, req *http.Request) {
	// var page int
	page, err := strconv.Atoi(req.FormValue("page"))

	if err == nil {
		pageSize, err := strconv.Atoi(req.FormValue("size"))

		if err == nil {
			hotelServiceRes, err := http.Get(
				fmt.Sprintf(
					"http://%s:%d/api/v1/hotels",
					service.reservationServiceHost,
					service.reservationServicePort,
				),
			)

			if err != nil {
				res.WriteHeader(http.StatusServiceUnavailable)
			} else {
				if hotelServiceRes.StatusCode == http.StatusOK {
					var resHotels []models.Hotel = parseHotelsRes(hotelServiceRes)
					var resultJson []byte = formHotelsResultJson(resHotels, page, pageSize)

					res.WriteHeader(http.StatusOK)
					res.Write(resultJson)
				} else {
					res.WriteHeader(hotelServiceRes.StatusCode)
				}
			}
		} else {
			res.WriteHeader(http.StatusBadRequest)
		}
	} else {
		res.WriteHeader(http.StatusBadRequest)
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
}

func (service *GatewayService) handleReservationsRequest(res http.ResponseWriter, req *http.Request) {
}

func (service *GatewayService) handleReservationUidRequest(res http.ResponseWriter, req *http.Request) {
}

func (service *GatewayService) handleLoyaltyRequest(res http.ResponseWriter, req *http.Request) {
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
