package services

import (
	"fmt"
	"net/http"

	"github.com/agarmirus/ds-lab02/internal/database"
	"github.com/agarmirus/ds-lab02/internal/models"
)

type ReservationService struct {
	hotelDAO  database.IDAO[models.Hotel]
	reservDAO database.IDAO[models.Reservation]

	port int
}

func NewReservationService(
	hotelDAO database.IDAO[models.Hotel],
	reservDAO database.IDAO[models.Reservation],
	port int,
) IService {
	return &ReservationService{hotelDAO, reservDAO, port}
}

func (service *ReservationService) Prepare() error {
	// TODO

	return nil
}

func (service *ReservationService) Run() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", service.port), nil)
}
