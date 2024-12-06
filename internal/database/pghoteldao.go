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

type PostgresHotelDAO struct {
	connStr string
}

func NewPostgresHotelDAO(connStr string) IDAO[models.Hotel] {
	return &PostgresHotelDAO{connStr}
}

func (dao *PostgresHotelDAO) SetConnectionString(connStr string) {
	dao.connStr = connStr
}

func (dao *PostgresHotelDAO) Create(hotel *models.Hotel) (models.Hotel, error) {
	log.Println("[ERROR] PostgresHotelDAO.Create. Method is not implemented")
	return models.Hotel{}, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresHotelDAO) Get() (resLst list.List, err error) {
	log.Println("[ERROR] PostgresHotelDAO.Get. Method is not implemented")
	return resLst, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresHotelDAO) GetPaginated(
	page int,
	pageSize int,
) (resLst list.List, err error) {
	if page <= 0 || pageSize <= 0 {
		log.Println("[ERROR] PostgresHotelDAO.GetPaginated. Invalid pages data")
		return resLst, errors.New(serverrors.ErrInvalidPagesData)
	}

	db, err := sql.Open(`postgres`, dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresHotelDAO.GetPaginated. Cannot connect to database:", err)
		return resLst, errors.New(serverrors.ErrDatabaseConnection)
	}

	rows, err := db.Query(
		`select * from hotels order by id limit $1 offset $2;`,
		pageSize,
		page*pageSize,
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("[ERROR] PostgresHotelDAO.GetPaginated. Error while executing query:", err)
		return resLst, errors.New(serverrors.ErrQueryResRead)
	}

	for rows.Next() {
		var hotel models.Hotel
		err = rows.Scan(&hotel)

		if err != nil {
			log.Println("[ERROR] PostgresHotelDAO.GetPaginated. Error while reading query result:", err)
			return list.List{}, errors.New(serverrors.ErrQueryResRead)
		}

		resLst.PushBack(hotel)
	}

	return resLst, nil
}

func (dao *PostgresHotelDAO) GetById(hotel *models.Hotel) (resHotel models.Hotel, err error) {
	if hotel.Id <= 0 {
		log.Println("[ERROR] PostgresHotelDAO.GetById. Invalid ID")
		return resHotel, errors.New(serverrors.ErrInvalidHotelId)
	}

	db, err := sql.Open(`postgres`, dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresHotelDAO.GetById. Cannot connect to database:", err)
		return resHotel, errors.New(serverrors.ErrDatabaseConnection)
	}

	row := db.QueryRow(
		`select * from hotels where id = $1;`,
		hotel.Id,
	)
	err = row.Scan(&resHotel)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[ERROR] PostgresHotelDAO.GetById. Entity not found")
			err = errors.New(serverrors.ErrEntityNotFound)
		} else {
			log.Println("[ERROR] PostgresHotelDAO.GetById. Error while reading query result:", err)
			err = errors.New(serverrors.ErrQueryResRead)
		}
	}

	return resHotel, err
}

func (dao *PostgresHotelDAO) GetByAttribute(attrName string, attrValue string) (resLst list.List, err error) {
	db, err := sql.Open(`postgres`, dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresHotelDAO.GetByAttribute. Cannot connect to database:", err)
		return resLst, errors.New(serverrors.ErrDatabaseConnection)
	}

	rows, err := db.Query(
		`select * from hotels where $1 = $2;`,
		attrName, attrValue,
	)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Println("[ERROR] PostgresHotelDAO.GetByAttribute. Error while executing query:", err)
			return resLst, errors.New(serverrors.ErrQueryResRead)
		}

		return resLst, nil
	}

	for rows.Next() {
		var hotel models.Hotel
		err = rows.Scan(&hotel)

		if err != nil {
			log.Println("[ERROR] PostgresHotelDAO.GetByAttribute. Error while reading query result:", err)
			return list.List{}, errors.New(serverrors.ErrQueryResRead)
		}

		resLst.PushBack(hotel)
	}

	return resLst, nil
}

func (dao *PostgresHotelDAO) Update(hotel *models.Hotel) (models.Hotel, error) {
	log.Println("[ERROR] PostgresHotelDAO.Update. Method is not implemented")
	return models.Hotel{}, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresHotelDAO) Delete(hotel *models.Hotel) error {
	log.Println("[ERROR] PostgresHotelDAO.Delete. Method is not implemented")
	return errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresHotelDAO) DeleteByAttr(attrName string, attrValue string) error {
	log.Println("[ERROR] PostgresHotelDAO.DeleteByAttr. Method is not implemented")
	return errors.New(serverrors.ErrMethodIsNotImplemented)
}
