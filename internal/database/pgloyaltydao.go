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

type PostgresLoyaltyDAO struct {
	connStr string
}

func NewPostgresLoyaltyDAO(connStr string) IDAO[models.Loyalty] {
	return &PostgresLoyaltyDAO{connStr}
}

func (dao *PostgresLoyaltyDAO) SetConnectionString(connStr string) {
	dao.connStr = connStr
}

func (dao *PostgresLoyaltyDAO) Create(loyalty *models.Loyalty) (models.Loyalty, error) {
	log.Println("[ERROR] PostgresLoyaltyDAO.Create. Method is not implemented")
	return models.Loyalty{}, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresLoyaltyDAO) Get() (list.List, error) {
	log.Println("[ERROR] PostgresLoyaltyDAO.Get. Method is not implemented")
	return list.List{}, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresLoyaltyDAO) GetById(loyalty *models.Loyalty) (models.Loyalty, error) {
	log.Println("[ERROR] PostgresLoyaltyDAO.GetById. Method is not implemented")
	return models.Loyalty{}, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresLoyaltyDAO) GetByAttribute(attrName string, attrValue string) (resLst list.List, err error) {
	db, err := sql.Open(`postgres`, dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresLoyaltyDAO.GetByAttribute. Cannot connect to database:", err)
		return resLst, errors.New(serverrors.ErrDatabaseConnection)
	}

	rows, err := db.Query(
		`select * from loyalty where $1 = $2;`,
		attrName, attrValue,
	)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Println("[ERROR] PostgresLoyaltyDAO.GetByAttribute. Error while executing query:", err)
			return resLst, errors.New(serverrors.ErrQueryResRead)
		}

		return resLst, nil
	}

	for rows.Next() {
		var loyalty models.Loyalty
		err = rows.Scan(&loyalty)

		if err != nil {
			log.Println("[ERROR] PostgresLoyaltyDAO.GetByAttribute. Error while reading query result:", err)
			return list.List{}, errors.New(serverrors.ErrQueryResRead)
		}

		resLst.PushBack(loyalty)
	}

	return resLst, nil
}

func (dao *PostgresLoyaltyDAO) Update(loyalty *models.Loyalty) (updatedLoyalty models.Loyalty, err error) {
	db, err := sql.Open(`postgres`, dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresLoyaltyDAO.Update. Cannot connect to database:", err)
		return updatedLoyalty, errors.New(serverrors.ErrDatabaseConnection)
	}

	row := db.QueryRow(
		`update loyalty
		set username = $1, reservation_count = $2, status = $3, discount = $4
		where id = $5
		returning *`,
		loyalty.Username, loyalty.ReservationCount, loyalty.Status, loyalty.Discount,
		loyalty.Id,
	)
	err = row.Scan(&updatedLoyalty)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[ERROR] PostgresLoyaltyDAO.Update. Entity not found")
			err = errors.New(serverrors.ErrEntityNotFound)
		} else {
			log.Println("[ERROR] PostgresLoyaltyDAO.Update. Error while reading query result:", err)
			err = errors.New(serverrors.ErrQueryResRead)
		}
	}

	return updatedLoyalty, err
}

func (dao *PostgresLoyaltyDAO) Delete(loyalty *models.Loyalty) error {
	log.Println("[ERROR] PostgresLoyaltyDAO.Delete. Method is not implemented")
	return errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresLoyaltyDAO) DeleteByAttr(attrName string, attrValue string) error {
	log.Println("[ERROR] PostgresLoyaltyDAO.DeleteByAttr. Method is not implemented")
	return errors.New(serverrors.ErrMethodIsNotImplemented)
}
