package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/agarmirus/ds-lab02/internal/models"
	"github.com/agarmirus/ds-lab02/internal/services"
	"github.com/google/uuid"
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

func (controller *GatewayController) handleAllHotelsGet(res http.ResponseWriter, req *http.Request) {
	page, pageParseErr := strconv.Atoi(req.FormValue(`page`))
	pageSize, pageSizeParseErr := strconv.Atoi(req.FormValue(`size`))

	if pageParseErr != nil || pageSizeParseErr != nil || page <= 0 || pageSize <= 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	pagRes, err := controller.service.ReadAllHotels(page, pageSize)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var pageResJSON []byte
	pageResJSON, err = json.Marshal(pagRes)

	if err == nil {
		res.WriteHeader(http.StatusOK)
		res.Header().Add(`Content-Type`, `application/json`)
		res.Write(pageResJSON)
	} else {
		res.WriteHeader(http.StatusInternalServerError)
	}
}

func (controller *GatewayController) handleUserInfoGet(res http.ResponseWriter, req *http.Request) {
	username := req.Header.Get(`X-User-Name`)

	if strings.Trim(username, ` `) == `` {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	userInfoRes, err := controller.service.ReadUserInfo(username)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var userInfoResJSON []byte
	userInfoResJSON, err = json.Marshal(userInfoRes)

	if err == nil {
		res.WriteHeader(http.StatusOK)
		res.Header().Add(`Content-Type`, `application/json`)
		res.Write(userInfoResJSON)
	} else {
		res.WriteHeader(http.StatusInternalServerError)
	}
}

func (controller *GatewayController) handleUserReservationsGet(res http.ResponseWriter, req *http.Request) {
	username := req.Header.Get(`X-User-Name`)

	if strings.Trim(username, ` `) == `` {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	reservsResSlice, err := controller.service.ReadUserReservations(username)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var reservsResSliceJSON []byte
	reservsResSliceJSON, err = json.Marshal(reservsResSlice)

	if err == nil {
		res.WriteHeader(http.StatusOK)
		res.Header().Add(`Content-Type`, `application/json`)
		res.Write(reservsResSliceJSON)
	} else {
		res.WriteHeader(http.StatusInternalServerError)
	}
}

func (controller *GatewayController) handleNewReservationPost(res http.ResponseWriter, req *http.Request) {
	username := req.Header.Get(`X-User-Name`)

	if strings.Trim(username, ` `) == `` {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var reqBody []byte
	n, err := req.Body.Read(reqBody)

	if err != nil || n <= 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var crReservReq models.CreateReservationRequest
	err = json.Unmarshal(reqBody, &crReservReq)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var validErrRes models.ValidationErrorResponse
	validErrRes, err = models.ValidateCrReservReq(&crReservReq)

	if err == nil {
		var crReservRes models.CreateReservationResponse
		crReservRes, err = controller.service.CreateReservation(username, &crReservReq)

		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		var crReservResJSON []byte
		crReservResJSON, err = json.Marshal(crReservRes)

		if err == nil {
			res.WriteHeader(http.StatusOK)
			res.Header().Add(`Content-Type`, `application/json`)
			res.Write(crReservResJSON)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		var validErrResJSON []byte
		validErrResJSON, err = json.Marshal(validErrRes)

		if err == nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Header().Add(`Content-Type`, `application/json`)
			res.Write(validErrResJSON)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (controller *GatewayController) handleSingleReservationGet(res http.ResponseWriter, req *http.Request) {
	reservationUid := req.PathValue("reservationUid")
	username := req.Header.Get(`X-User-Name`)

	if strings.Trim(username, ` `) == `` || uuid.Validate(reservationUid) != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	reservRes, err := controller.service.ReadReservation(reservationUid, username)

	if err != nil {
		// TODO: Check for no results case
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	reservResJSON, err := json.Marshal(reservRes)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Header().Add(`Content-Type`, `application/json`)
	res.Write(reservResJSON)
}

func (controller *GatewayController) handleSingleReservationDelete(res http.ResponseWriter, req *http.Request) {
	reservationUid := req.PathValue("reservationUid")
	username := req.Header.Get(`X-User-Name`)

	if strings.Trim(username, ` `) == `` || uuid.Validate(reservationUid) != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err := controller.service.DeleteReservation(reservationUid, username)

	if err != nil {
		// TODO: Check for no results case
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}

func (controller *GatewayController) handleLoyaltyGet(res http.ResponseWriter, req *http.Request) {
	username := req.Header.Get(`X-User-Name`)

	if strings.Trim(username, ` `) == `` {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	loyaltyInfoRes, err := controller.service.ReadUserLoyalty(username)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	loyaltyInfoResJSON, err := json.Marshal(loyaltyInfoRes)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Header().Add(`Content-Type`, `application/json`)
	res.Write(loyaltyInfoResJSON)
}

func (controller *GatewayController) handleHotelsRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		controller.handleAllHotelsGet(res, req)
	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *GatewayController) handleUserRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		controller.handleUserInfoGet(res, req)
	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *GatewayController) handleReservationsRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		controller.handleUserReservationsGet(res, req)
	} else if req.Method == `POST` {
		controller.handleNewReservationPost(res, req)
	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *GatewayController) handleSingleReservationRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		controller.handleSingleReservationGet(res, req)
	} else if req.Method == `DELETE` {
		controller.handleSingleReservationDelete(res, req)
	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *GatewayController) handleLoyaltyRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		controller.handleLoyaltyGet(res, req)
	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
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
