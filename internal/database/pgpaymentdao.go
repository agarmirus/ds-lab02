package database

import (
	"container/list"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx"

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
	db, err := sql.Open(`postgres`, dao.connStr)

	if err != nil {
		return resLst, errors.New(serverrors.ErrDatabaseConnection)
	}

	rows, err := db.Query(
		`select * from payment where $1 = '$2';`,
		attrName, attrValue,
	)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return resLst, errors.New(serverrors.ErrQueryResRead)
		}

		return resLst, nil
	}

	for rows.Next() {
		var payment models.Payment
		err = rows.Scan(&payment)

		if err != nil {
			return list.List{}, errors.New(serverrors.ErrQueryResRead)
		}

		resLst.PushBack(payment)
	}

	return resLst, nil
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
