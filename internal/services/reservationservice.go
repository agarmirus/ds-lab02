package services

import (
	"fmt"
	"net/http"

	"github.com/agarmirus/ds-lab02/internal/database"
	"github.com/agarmirus/ds-lab02/internal/models"
)

type HotelDAO = database.IDAO[models.Hotel]
type ReservationDAO = database.IDAO[models.Reservation]

type ReservationService struct {
	hotelDAO  HotelDAO
	reservDAO ReservationDAO

	port int
}

func NewReservationService(
	hotelDAO HotelDAO,
	reservDAO ReservationDAO,
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
