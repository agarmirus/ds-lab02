package services

import (
	"fmt"
	"net/http"

	"github.com/agarmirus/ds-lab02/internal/database"
	"github.com/agarmirus/ds-lab02/internal/models"
)

type LoyaltyService struct {
	loyaltyService database.IDAO[models.Loyalty]

	port int
}

func NewLoyaltyService(
	loyaltyService database.IDAO[models.Loyalty],
	port int,
) IService {
	return &LoyaltyService{loyaltyService, port}
}

func (service *LoyaltyService) Prepare() error {
	// TODO

	return nil
}

func (service *LoyaltyService) Run() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", service.port), nil)
}
