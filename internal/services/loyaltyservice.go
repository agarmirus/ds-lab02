package services

import (
	"errors"

	"github.com/agarmirus/ds-lab02/internal/database"
	"github.com/agarmirus/ds-lab02/internal/models"
	"github.com/agarmirus/ds-lab02/internal/serverrors"
)

type LoyaltyService struct {
	loyaltyDAO database.IDAO[models.Loyalty]
}

func NewLoyaltyService(
	loyaltyDAO database.IDAO[models.Loyalty],
) ILoyaltyService {
	return &LoyaltyService{loyaltyDAO}
}

func (service *LoyaltyService) ReadLoyaltyByUsername(username string) (loyalty models.Loyalty, err error) {
	loyaltiesLst, err := service.loyaltyDAO.GetByAttribute(`username`, username)

	if err != nil {
		return loyalty, err
	}

	if loyaltiesLst.Len() == 0 {
		return loyalty, errors.New(serverrors.ErrEntityNotFound)
	}

	return loyaltiesLst.Front().Value.(models.Loyalty), nil
}

func (service *LoyaltyService) UpdateLoyaltyById(loyalty *models.Loyalty) (updatedLoyalty models.Loyalty, err error) {
	return service.loyaltyDAO.Update(loyalty)
}