package services

import (
	"errors"

	"github.com/agarmirus/ds-lab02/internal/database"
	"github.com/agarmirus/ds-lab02/internal/models"
	"github.com/agarmirus/ds-lab02/internal/serverrors"
)

type PaymentService struct {
	paymentDAO database.IDAO[models.Payment]
}

func NewPaymentService(
	paymentDAO database.IDAO[models.Payment],
) IPaymentService {
	return &PaymentService{paymentDAO}
}

func (service *PaymentService) ReadPaymentByUid(paymentUid string) (payment models.Payment, err error) {
	paymentsLst, err := service.paymentDAO.GetByAttribute(`payment_uid`, paymentUid)

	if err != nil {
		return payment, err
	}

	if paymentsLst.Len() == 0 {
		return payment, errors.New(serverrors.ErrEntityNotFound)
	}

	return paymentsLst.Front().Value.(models.Payment), nil
}

func (service *PaymentService) UpdatePaymentByUid(payment *models.Payment) (newPayment models.Payment, err error) {
	return service.paymentDAO.Update(payment)
}

func (service *PaymentService) CreatePayment(payment *models.Payment) (newPayment models.Payment, err error) {
	return service.paymentDAO.Create(payment)
}
