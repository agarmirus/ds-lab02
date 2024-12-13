package database

import (
	"container/list"
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"

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
	return models.Hotel{}, serverrors.ErrMethodIsNotImplemented
}

func (dao *PostgresHotelDAO) Get() (resLst list.List, err error) {
	log.Println("[ERROR] PostgresHotelDAO.Get. Method is not implemented")
	return resLst, serverrors.ErrMethodIsNotImplemented
}

func (dao *PostgresHotelDAO) GetPaginated(
	page int,
	pageSize int,
) (resLst list.List, err error) {
	if page <= 0 || pageSize <= 0 {
		log.Println("[ERROR] PostgresHotelDAO.GetPaginated. Invalid pages data")
		return resLst, serverrors.ErrInvalidPagesData
	}

	db, err := sql.Open(`postgres`, dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresHotelDAO.GetPaginated. Cannot connect to database:", err)
		return resLst, serverrors.ErrDatabaseConnection
	}

	defer db.Close()

	rows, err := db.Query(
		`select * from hotels order by id limit $1 offset $2;`,
		pageSize,
		(page-1)*pageSize,
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("[ERROR] PostgresHotelDAO.GetPaginated. Error while executing query:", err)
		return resLst, serverrors.ErrQueryResRead
	}

	defer rows.Close()

	for rows.Next() {
		var hotel models.Hotel
		err = rows.Scan(
			&hotel.Id, &hotel.Uid,
			&hotel.Name, &hotel.Country,
			&hotel.City, &hotel.Address,
			&hotel.Stars, &hotel.Price,
		)

		if err != nil {
			log.Println("[ERROR] PostgresHotelDAO.GetPaginated. Error while reading query result:", err)
			return list.List{}, serverrors.ErrQueryResRead
		}

		resLst.PushBack(hotel)
	}

	return resLst, nil
}

func (dao *PostgresHotelDAO) GetById(hotel *models.Hotel) (resHotel models.Hotel, err error) {
	if hotel.Id <= 0 {
		log.Println("[ERROR] PostgresHotelDAO.GetById. Invalid ID")
		return resHotel, serverrors.ErrInvalidHotelId
	}

	db, err := sql.Open(`postgres`, dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresHotelDAO.GetById. Cannot connect to database:", err)
		return resHotel, serverrors.ErrDatabaseConnection
	}

	defer db.Close()

	row := db.QueryRow(
		`select * from hotels where id = $1;`,
		hotel.Id,
	)

	err = row.Scan(
		&resHotel.Id, &resHotel.Uid,
		&resHotel.Name, &resHotel.Country,
		&resHotel.City, &resHotel.Address,
		&resHotel.Stars, &resHotel.Price,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[ERROR] PostgresHotelDAO.GetById. Entity not found")
			err = serverrors.ErrEntityNotFound
		} else {
			log.Println("[ERROR] PostgresHotelDAO.GetById. Error while reading query result:", err)
			err = serverrors.ErrQueryResRead
		}
	}

	return resHotel, err
}

func (dao *PostgresHotelDAO) GetByAttribute(attrName string, attrValue string) (resLst list.List, err error) {
	db, err := sql.Open(`postgres`, dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresHotelDAO.GetByAttribute. Cannot connect to database:", err)
		return resLst, serverrors.ErrDatabaseConnection
	}

	defer db.Close()

	queryStr := fmt.Sprintf(
		`select * from hotels where %s = $1;`,
		attrName,
	)
	rows, err := db.Query(queryStr, attrValue)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Println("[ERROR] PostgresHotelDAO.GetByAttribute. Error while executing query:", err)
			return resLst, serverrors.ErrQueryResRead
		}

		return resLst, nil
	}

	defer rows.Close()

	for rows.Next() {
		var hotel models.Hotel
		err = rows.Scan(
			&hotel.Id, &hotel.Uid,
			&hotel.Name, &hotel.Country,
			&hotel.City, &hotel.Address,
			&hotel.Stars, &hotel.Price,
		)

		if err != nil {
			log.Println("[ERROR] PostgresHotelDAO.GetByAttribute. Error while reading query result:", err)
			return list.List{}, serverrors.ErrQueryResRead
		}

		resLst.PushBack(hotel)
	}

	return resLst, nil
}

func (dao *PostgresHotelDAO) Update(hotel *models.Hotel) (models.Hotel, error) {
	log.Println("[ERROR] PostgresHotelDAO.Update. Method is not implemented")
	return models.Hotel{}, serverrors.ErrMethodIsNotImplemented
}

func (dao *PostgresHotelDAO) Delete(hotel *models.Hotel) error {
	log.Println("[ERROR] PostgresHotelDAO.Delete. Method is not implemented")
	return serverrors.ErrMethodIsNotImplemented
}

func (dao *PostgresHotelDAO) DeleteByAttr(attrName string, attrValue string) error {
	log.Println("[ERROR] PostgresHotelDAO.DeleteByAttr. Method is not implemented")
	return serverrors.ErrMethodIsNotImplemented
}
