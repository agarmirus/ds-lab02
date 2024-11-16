package models

import (
	"errors"
	"time"

	"github.com/agarmirus/ds-lab02/internal/serverrors"
	"github.com/google/uuid"
)

type CreateReservationRequest struct {
	HotelUid  string `json:"hotelUid"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

func ValidateCrReservReq(
	createReservReq *CreateReservationRequest,
) (validErrRes ValidationErrorResponse, err error) {
	if uuid.Validate(createReservReq.HotelUid) != nil {
		validErrRes.Errors = append(validErrRes.Errors, ErrorDiscription{Field: `hotelUid`, Error: `invalid uid`})
	}

	startDate, err := time.Parse(`%F`, createReservReq.StartDate)

	if err != nil {
		validErrRes.Errors = append(validErrRes.Errors, ErrorDiscription{Field: `startDate`, Error: `invalid date format`})
	}

	endDate, err := time.Parse(`%F`, createReservReq.EndDate)

	if err != nil {
		validErrRes.Errors = append(validErrRes.Errors, ErrorDiscription{Field: `endDate`, Error: `invalid date format`})
	} else if startDate.Unix() > endDate.Unix() {
		validErrRes.Errors = append(validErrRes.Errors, ErrorDiscription{Field: `startDate`, Error: `invalid date period`})
		err = errors.New(serverrors.ErrInvalidReservDates)
	}

	if err != nil {
		validErrRes.Message = `invalid reservation request data`
	}

	return validErrRes, err
}
