package services

import (
	"container/list"
	"errors"

	"github.com/agarmirus/ds-lab02/internal/database"
	"github.com/agarmirus/ds-lab02/internal/models"
	"github.com/agarmirus/ds-lab02/internal/serverrors"
)

type ReservationService struct {
	reservsDAO database.IDAO[models.Reservation]
	hotelsDAO  database.IDAO[models.Hotel]
}

func NewReservationService(
	reservsDAO database.IDAO[models.Reservation],
	hotelsDAO database.IDAO[models.Hotel],
) IReservationService {
	return &ReservationService{reservsDAO, hotelsDAO}
}

func (service *ReservationService) ReadAllHotels() (hotelsLst list.List, err error) {
	return service.hotelsDAO.Get()
}

func (service *ReservationService) ReadHotelById(hotelId int) (hotel models.Hotel, err error) {
	desiredHotel := models.Hotel{Id: hotelId}
	return service.hotelsDAO.GetById(&desiredHotel)
}

func (service *ReservationService) ReadHotelByUid(hotelUid string) (hotel models.Hotel, err error) {
	hotelsLst, err := service.hotelsDAO.GetByAttribute(`hotel_uid`, hotelUid)

	if err != nil {
		return hotel, err
	}

	if hotelsLst.Len() == 0 {
		return hotel, errors.New(serverrors.ErrEntityNotFound)
	}

	return hotelsLst.Front().Value.(models.Hotel), nil
}

func (service *ReservationService) ReadReservsByUsername(username string) (reservsLst list.List, err error) {
	reservsLst, err = service.hotelsDAO.GetByAttribute(`username`, username)

	if err != nil {
		return reservsLst, err
	}

	if reservsLst.Len() == 0 {
		return reservsLst, errors.New(serverrors.ErrEntityNotFound)
	}

	return reservsLst, nil
}

func (service *ReservationService) ReadReservByUid(reservUid string) (reservation models.Reservation, err error) {
	reservsLst, err := service.hotelsDAO.GetByAttribute(`reservation_uid`, reservUid)

	if err != nil {
		return reservation, err
	}

	if reservsLst.Len() == 0 {
		return reservation, errors.New(serverrors.ErrEntityNotFound)
	}

	return reservsLst.Front().Value.(models.Reservation), nil
}

func (service *ReservationService) UpdateReservByUid(reservation *models.Reservation) (updatedReservation models.Reservation, err error) {
	return service.reservsDAO.Update(reservation)
}

func (service *ReservationService) CreateReserv(reservation *models.Reservation) (newReservation models.Reservation, err error) {
	return service.reservsDAO.Create(reservation)
}
