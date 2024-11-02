package services

import (
	"fmt"
	"net/http"

	"github.com/agarmirus/ds-lab02/internal/database"
	"github.com/agarmirus/ds-lab02/internal/models"
)

type PaymentService struct {
	paymentDAO database.IDAO[models.Payment]

	port int
}

func NewPaymentService(
	paymentDAO database.IDAO[models.Payment],
	port int,
) IService {
	return &PaymentService{paymentDAO, port}
}

func (service *PaymentService) Prepare() error {
	// TODO

	return nil
}

func (service *PaymentService) Run() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", service.port), nil)
}
