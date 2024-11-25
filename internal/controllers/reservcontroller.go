package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/agarmirus/ds-lab02/internal/models"
	"github.com/agarmirus/ds-lab02/internal/serverrors"
	"github.com/agarmirus/ds-lab02/internal/services"
)

type ReservationController struct {
	host string
	port int

	service services.IReservationService
}

func NewReservationController(
	host string,
	port int,
	service services.IGatewayService,
) IController {
	return &ReservationController{host, port, service}
}

func (controller *ReservationController) handleAllHotelsGet(res http.ResponseWriter, req *http.Request) {
	hotelsLst, err := controller.service.ReadAllHotels()

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var hotelsSlice []models.Hotel
	hotelsLstEl := hotelsLst.Front()

	for hotelsLstEl != nil {
		hotelsSlice = append(hotelsSlice, hotelsLstEl.Value.(models.Hotel))
		hotelsLstEl.Next()
	}

	hotelsSliceJSON, err := json.Marshal(hotelsSlice)

	res.WriteHeader(http.StatusOK)
	res.Header().Add(`Content-Type`, `application/json`)
	res.Write(hotelsSliceJSON)
}

func (controller *ReservationController) handleHotelByIdGet(res http.ResponseWriter, req *http.Request) {
	hotelId, err := strconv.Atoi(req.Header.Get(`Hotel-Id`))

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	hotel, err := controller.service.ReadHotelById(hotelId)

	if err != nil {
		if errors.Is(err, errors.New(serverrors.ErrEntityNotFound)) {
			res.WriteHeader(http.StatusNotFound)
			return
		} else {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	hotelJSON, err := json.Marshal(hotel)

	res.WriteHeader(http.StatusOK)
	res.Header().Add(`Content-Type`, `application/json`)
	res.Write(hotelJSON)
}

func (controller *ReservationController) hangleHotelByUidGet(res http.ResponseWriter, req *http.Request) {
	hotelUid := req.PathValue("hotelUid")

	if hotelUid == `` {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	hotel, err := controller.service.ReadHotelByUid(hotelUid)

	if err != nil {
		if errors.Is(err, errors.New(serverrors.ErrEntityNotFound)) {
			res.WriteHeader(http.StatusNotFound)
			return
		} else {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	hotelJSON, err := json.Marshal(hotel)

	res.WriteHeader(http.StatusOK)
	res.Header().Add(`Content-Type`, `application/json`)
	res.Write(hotelJSON)
}

func (controller *ReservationController) handleReservsByUsernameGet(res http.ResponseWriter, req *http.Request) {
	reservUid := req.PathValue("reservUid")

	if reservUid == `` {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	reservation, err := controller.service.ReadReservByUid(reservUid)

	if err != nil {
		if errors.Is(err, errors.New(serverrors.ErrEntityNotFound)) {
			res.WriteHeader(http.StatusNotFound)
			return
		} else {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	reservJSON, err := json.Marshal(reservation)

	res.WriteHeader(http.StatusOK)
	res.Header().Add(`Content-Type`, `application/json`)
	res.Write(reservJSON)
}

func (controller *ReservationController) handleReservByUidGet(res http.ResponseWriter, req *http.Request) {
	username := req.Header.Get(`X-User-Name`)

	if username == `` {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	reservsLst, err := controller.service.ReadReservsByUsername(username)

	if err != nil {
		if errors.Is(err, errors.New(serverrors.ErrEntityNotFound)) {
			res.WriteHeader(http.StatusNotFound)
			return
		} else {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	var reservsSlice []models.Reservation
	reservsLstEl := reservsLst.Front()

	for reservsLstEl != nil {
		reservsSlice = append(reservsSlice, reservsLstEl.Value.(models.Reservation))
		reservsLstEl.Next()
	}

	reservsSliceJSON, err := json.Marshal(reservsSlice)

	res.WriteHeader(http.StatusOK)
	res.Header().Add(`Content-Type`, `application/json`)
	res.Write(reservsSliceJSON)
}

func (controller *ReservationController) handleHotelsRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		if res.Header().Get(`Hotel-Id`) == `` {
			controller.handleAllHotelsGet(res, req)
		} else {
			controller.handleHotelByIdGet(res, req)
		}
	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *ReservationController) handleHotelWithUidRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		controller.hangleHotelByUidGet(res, req)
	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *ReservationController) handleReservsRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		if res.Header().Get(`X-User-Name`) != `` {
			controller.handleReservsByUsernameGet(res, req)
		} else {
			res.WriteHeader(http.StatusMethodNotAllowed)
		}
	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *ReservationController) handleReservWithUidRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		controller.handleReservByUidGet(res, req)
	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *ReservationController) Prepare() error {
	http.HandleFunc(`/api/v1/hotels`, controller.handleHotelsRequest)
	http.HandleFunc(`/api/v1/hotels/{hotelUid}`, controller.handleHotelWithUidRequest)
	http.HandleFunc(`/api/v1/reservations`, controller.handleReservsRequest)
	http.HandleFunc(`/api/v1/reservations/{reservUid}`, controller.handleReservWithUidRequest)

	return nil
}

func (controller *ReservationController) Run() error {
	return http.ListenAndServe(
		fmt.Sprintf(
			`%s:%d`,
			controller.host,
			controller.port,
		),
		nil,
	)
}
