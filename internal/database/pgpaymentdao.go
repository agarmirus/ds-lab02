package database

import (
	"container/list"
	"database/sql"
	"errors"

	"github.com/agarmirus/ds-lab02/internal/models"
	"github.com/agarmirus/ds-lab02/internal/serverrors"
)

type PostgresPaymentDAO struct {
	connStr string
}

func NewPostgresPaymentDAO(connStr string) IDAO[models.Payment] {
	return &PostgresPaymentDAO{connStr}
}

func (dao *PostgresPaymentDAO) SetConnectionString(connStr string) {
	dao.connStr = connStr
}

func (dao *PostgresPaymentDAO) Create(payment *models.Payment) (models.Payment, error) {
	return models.Payment{}, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresPaymentDAO) Get() (list.List, error) {
	return list.List{}, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresPaymentDAO) GetById(payment *models.Payment) (models.Payment, error) {
	return models.Payment{}, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresPaymentDAO) GetByAttribute(attrName string, attrValue string) (resLst list.List, err error) {
	db, localErr := sql.Open(`postgres`, dao.connStr)

	if localErr == nil {
		var rows *sql.Rows
		rows, localErr = db.Query(
			`select * from payment where $1 = '$2';`,
			attrName, attrValue,
		)
		for localErr == nil && rows.Next() {
			var payment models.Payment
			localErr = rows.Scan(&payment)

			if localErr != nil {
				err = errors.New(serverrors.ErrQueryResRead)
			}

			resLst.PushBack(payment)
		}
	} else {
		err = errors.New(serverrors.ErrDatabaseConnection)
	}

	return resLst, err
}

func (dao *PostgresPaymentDAO) Update(payment *models.Payment) (models.Payment, error) {
	return models.Payment{}, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresPaymentDAO) Delete(payment *models.Payment) error {
	return errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresPaymentDAO) DeleteByAttr(attrName string, attrValue string) error {
	return errors.New(serverrors.ErrMethodIsNotImplemented)
}
