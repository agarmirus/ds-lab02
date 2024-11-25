package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/agarmirus/ds-lab02/internal/models"
	"github.com/agarmirus/ds-lab02/internal/serverrors"
	"github.com/agarmirus/ds-lab02/internal/services"
	"github.com/google/uuid"
)

type PaymentController struct {
	host string
	port int

	service services.IPaymentService
}

func NewPaymentController(
	host string,
	port int,
	service services.IPaymentService,
) IController {
	return &PaymentController{host, port, service}
}

func (controller *PaymentController) handlePaymentByPricePost(res http.ResponseWriter, req *http.Request) {
	price, err := strconv.Atoi(req.Header.Get(`Price`))

	if err != nil || price <= 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	payment := models.Payment{Uid: uuid.New().String(), Status: `PAID`, Price: price}

	newPayment, err := controller.service.CreatePayment(&payment)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}

	newPaymentJSON, err := json.Marshal(newPayment)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}

	res.WriteHeader(http.StatusOK)
	res.Header().Add(`Content-Type`, `application/json`)
	res.Write(newPaymentJSON)
}

func (controller *PaymentController) handlePaymentByUidGet(res http.ResponseWriter, req *http.Request) {
	paymentUid := req.PathValue("paymentUid")

	if paymentUid == `` {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	payment, err := controller.service.ReadPaymentByUid(paymentUid)

	if err != nil {
		if errors.Is(err, errors.New(serverrors.ErrEntityNotFound)) {
			res.WriteHeader(http.StatusNotFound)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	paymentJSON, err := json.Marshal(payment)

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}

	res.WriteHeader(http.StatusOK)
	res.Header().Add(`Content-Type`, `application/json`)
	res.Write(paymentJSON)
}

func (controller *PaymentController) handlePaymentByUidPatch(res http.ResponseWriter, req *http.Request) {
	paymentUid := req.PathValue("paymentUid")

	if strings.Trim(paymentUid, ` `) == `` {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var reqBody []byte
	n, err := req.Body.Read(reqBody)

	if err != nil || n <= 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var payment models.Payment
	err = json.Unmarshal(reqBody, &payment)

	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	payment.Uid = paymentUid
	_, err = controller.service.UpdatePaymentByUid(&payment)

	if err != nil {
		if errors.Is(err, errors.New(serverrors.ErrEntityNotFound)) {
			res.WriteHeader(http.StatusNotFound)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	res.WriteHeader(http.StatusOK)
}

func (controller *PaymentController) handlePaymentRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `POST` {
		if strings.Trim(req.Header.Get(`Price`), ` `) != `` {
			controller.handlePaymentByPricePost(res, req)
		} else {
			res.WriteHeader(http.StatusBadRequest)
		}
	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *PaymentController) handlePaymentByUidRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		controller.handlePaymentByUidGet(res, req)
	} else if req.Method == `PATCH` {
		controller.handlePaymentByUidPatch(res, req)
	} else {
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *PaymentController) Prepare() error {
	http.HandleFunc(`/api/v1/payment`, controller.handlePaymentRequest)
	http.HandleFunc(`/api/v1/payment/{paymentUid}`, controller.handlePaymentByUidRequest)

	return nil
}

func (controller *PaymentController) Run() error {
	return http.ListenAndServe(
		fmt.Sprintf(
			`%s:%d`,
			controller.host,
			controller.port,
		),
		nil,
	)
}
