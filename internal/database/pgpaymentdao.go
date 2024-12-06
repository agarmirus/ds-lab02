package database

import (
	"container/list"
	"database/sql"
	"errors"
	"log"

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
	log.Println("[ERROR] PostgresPaymentDAO.Create. Method is not implemented")
	return models.Payment{}, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresPaymentDAO) Get() (list.List, error) {
	log.Println("[ERROR] PostgresPaymentDAO.Get. Method is not implemented")
	return list.List{}, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresPaymentDAO) GetPaginated(
	page int,
	pageSize int,
) (resLst list.List, err error) {
	log.Println("[ERROR] PostgresLoyaltyDAO.GetPaginated. Method is not implemented")
	return list.List{}, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresPaymentDAO) GetById(payment *models.Payment) (models.Payment, error) {
	log.Println("[ERROR] PostgresPaymentDAO.GetById. Method is not implemented")
	return models.Payment{}, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresPaymentDAO) GetByAttribute(attrName string, attrValue string) (resLst list.List, err error) {
	db, err := sql.Open(`postgres`, dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresPaymentDAO.GetByAttribute. Cannot connect to database:", err)
		return resLst, errors.New(serverrors.ErrDatabaseConnection)
	}

	rows, err := db.Query(
		`select * from payment where $1 = $2;`,
		attrName, attrValue,
	)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Println("[ERROR] PostgresPaymentDAO.GetByAttribute. Error while executing query:", err)
			return resLst, errors.New(serverrors.ErrQueryResRead)
		}

		return resLst, nil
	}

	for rows.Next() {
		var payment models.Payment
		err = rows.Scan(&payment)

		if err != nil {
			log.Println("[ERROR] PostgresPaymentDAO.GetByAttribute. Error while reading query result:", err)
			return list.List{}, errors.New(serverrors.ErrQueryResRead)
		}

		resLst.PushBack(payment)
	}

	return resLst, nil
}

func (dao *PostgresPaymentDAO) Update(payment *models.Payment) (updatedPayment models.Payment, err error) {
	db, err := sql.Open(`postgres`, dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresPaymentDAO.Update. Cannot connect to database:", err)
		return updatedPayment, errors.New(serverrors.ErrDatabaseConnection)
	}

	row := db.QueryRow(
		`update payment
		set status = $1, price = $2
		where payment_uid = $3
		returning *`,
		payment.Status, payment.Price,
		payment.Uid,
	)
	err = row.Scan(&updatedPayment)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[ERROR] PostgresPaymentDAO.Update. Entity not found")
			err = errors.New(serverrors.ErrEntityNotFound)
		} else {
			log.Println("[ERROR] PostgresPaymentDAO.Update. Error while reading query result:", err)
			err = errors.New(serverrors.ErrQueryResRead)
		}
	}

	return updatedPayment, err
}

func (dao *PostgresPaymentDAO) Delete(payment *models.Payment) error {
	log.Println("[ERROR] PostgresPaymentDAO.Delete. Method is not implemented")
	return errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresPaymentDAO) DeleteByAttr(attrName string, attrValue string) error {
	log.Println("[ERROR] PostgresPaymentDAO.DeleteByAttr. Method is not implemented")
	return errors.New(serverrors.ErrMethodIsNotImplemented)
}
