package database

import (
	"container/list"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/agarmirus/ds-lab02/internal/models"
)

type PostgresPaymentDAO struct {
	connStr string
}

func NewPostgresPaymentDAO() IDAO[models.Payment] {
	return &PostgresPaymentDAO{}
}

func (dao *PostgresPaymentDAO) SetConnectionString(connStr string) {
	dao.connStr = connStr
}

func (dao *PostgresPaymentDAO) Create(payment *models.Payment) error {
	var err error = nil

	db, err := sql.Open("postgres", dao.connStr)

	if err == nil {
		defer db.Close()

		queryString := fmt.Sprintf(
			"insert into payment (payment_uid, status, price) values ('%s', '%s', %d);",
			payment.GetUid().String(),
			payment.GetStatus(),
			payment.GetPrice(),
		)

		_, err = db.Exec(queryString)
	}

	return err
}

func (dao *PostgresPaymentDAO) Get() (list.List, error) {
	return list.List{}, errors.New("PostgresPaymentDAO.Get() method is not implemented")
}

func (dao *PostgresPaymentDAO) GetById(payment *models.Payment) (models.Payment, error) {
	var err error = nil
	result := models.Payment{}

	db, err := sql.Open("postgres", dao.connStr)

	if err == nil {
		defer db.Close()

		queryString := fmt.Sprintf(
			"select * from payment where id = %d;",
			payment.GetId(),
		)

		var row *sql.Row = db.QueryRow(queryString)

		var id, price int
		var uid uuid.UUID
		var status string

		err = row.Scan(&id, &uid, &status, &price)

		if err == nil {
			result.SetId(id)
			result.SetUid(uid)
			result.SetStatus(status)
			result.SetPrice(price)
		} else if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
	}

	return result, err
}

func (dao *PostgresPaymentDAO) GetByAttribute(attrName string, attrValue string) (list.List, error) {
	var err error = nil
	resLst := list.List{}

	db, err := sql.Open("postgres", dao.connStr)

	if err == nil {
		defer db.Close()

		queryString := fmt.Sprintf(
			"select * from payment where %s = '%s';",
			attrName,
			attrValue,
		)

		var rows *sql.Rows
		rows, err = db.Query(queryString)

		if err == nil {
			defer rows.Close()

			for rows.Next() {
				var id, price int
				var uid uuid.UUID
				var status string

				err = rows.Scan(&id, &uid, &status, &price)

				if err == nil {
					curPayment := models.Payment{}
					curPayment.SetId(id)
					curPayment.SetUid(uid)
					curPayment.SetStatus(status)
					curPayment.SetPrice(price)

					resLst.PushBack(curPayment)
				} else if errors.Is(err, sql.ErrNoRows) {
					err = nil
				}
			}
		}
	}

	return resLst, err
}

func (dao *PostgresPaymentDAO) Update(*models.Payment) (models.Payment, error) {
	return models.Payment{}, errors.New("PostgresPaymentDAO.Update() method is not implemented")
}

func (dao *PostgresPaymentDAO) Delete(payment *models.Payment) error {
	return errors.New("PostgresPaymentDAO.Delete() method is not implemented")
}
