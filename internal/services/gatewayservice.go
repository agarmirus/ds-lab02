package services

import (
	"fmt"
	"net/http"
)

type GatewayService struct {
	port int
}

func NewGatewayService(port int) IService {
	return &GatewayService{port}
}

func (service *GatewayService) handleHotelsRequest(res http.ResponseWriter, req *http.Request) {
}

func (service *GatewayService) handleUserRequest(res http.ResponseWriter, req *http.Request) {
}

func (service *GatewayService) handleReservationsRequest(res http.ResponseWriter, req *http.Request) {
}

func (service *GatewayService) handleLoyaltyRequest(res http.ResponseWriter, req *http.Request) {
}

func (service *GatewayService) Prepare() error {
	http.HandleFunc("/api/v1/hotels", service.handleHotelsRequest)
	http.HandleFunc("/api/v1/me", service.handleUserRequest)
	http.HandleFunc("/api/v1/reservations", service.handleReservationsRequest)
	http.HandleFunc("/api/v1/loyalty", service.handleLoyaltyRequest)

	return nil
}

func (service *GatewayService) Run() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", service.port), nil)
}
