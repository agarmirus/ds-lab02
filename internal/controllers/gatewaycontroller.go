package controllers

import (
	"fmt"
	"net/http"

	"github.com/agarmirus/ds-lab02/internal/services"
)

type GatewayController struct {
	host string
	port int

	service services.IGatewayService
}

func NewGatewayController(
	host string,
	port int,
	service services.IGatewayService,
) IController {
	return &GatewayController{host, port, service}
}

func (controller *GatewayController) handleHotelsRequest(res http.ResponseWriter, req *http.Request) {
}

func (controller *GatewayController) handleUserRequest(res http.ResponseWriter, req *http.Request) {
}

func (controller *GatewayController) handleReservationsRequest(res http.ResponseWriter, req *http.Request) {
}

func (controller *GatewayController) handleSingleReservationRequest(res http.ResponseWriter, req *http.Request) {
}

func (controller *GatewayController) handleLoyaltyRequest(res http.ResponseWriter, req *http.Request) {
}

func (controller *GatewayController) Prepare() error {
	http.HandleFunc(`/api/v1/hotels`, controller.handleHotelsRequest)
	http.HandleFunc(`/api/v1/me`, controller.handleUserRequest)
	http.HandleFunc(`/api/v1/reservations`, controller.handleReservationsRequest)
	http.HandleFunc(`/api/v1/reservations/{reservationUid}`, controller.handleSingleReservationRequest)
	http.HandleFunc(`/api/v1/loyalty`, controller.handleLoyaltyRequest)

	return nil
}

func (controller *GatewayController) Run() error {
	return http.ListenAndServe(
		fmt.Sprintf(
			`%s:%d`,
			controller.host,
			controller.port,
		),
		nil,
	)
}
