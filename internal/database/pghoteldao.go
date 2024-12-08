package database

import (
	"container/list"
	"database/sql"
	"errors"

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
	return models.Hotel{}, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresHotelDAO) Get() (resLst list.List, err error) {
	db, err := sql.Open(`postgres`, dao.connStr)

	if err != nil {
		return resLst, errors.New(serverrors.ErrDatabaseConnection)
	}

	rows, err := db.Query(`select * from hotels;`)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return resLst, errors.New(serverrors.ErrQueryResRead)
	}

	for rows.Next() {
		var hotel models.Hotel
		err = rows.Scan(&hotel)

		if err != nil {
			return list.List{}, errors.New(serverrors.ErrQueryResRead)
		}

		resLst.PushBack(hotel)
	}

	return resLst, nil
}

func (dao *PostgresHotelDAO) GetById(hotel *models.Hotel) (resHotel models.Hotel, err error) {
	if hotel.Id <= 0 {
		return resHotel, errors.New(serverrors.ErrInvalidHotelId)
	}

	db, err := sql.Open(`postgres`, dao.connStr)

	if err != nil {
		return resHotel, errors.New(serverrors.ErrDatabaseConnection)
	}

	row := db.QueryRow(
		`select * from hotels where id = $1;`,
		hotel.Id,
	)
	err = row.Scan(&resHotel)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.New(serverrors.ErrEntityNotFound)
		} else {
			err = errors.New(serverrors.ErrQueryResRead)
		}
	}

	return resHotel, err
}

func (dao *PostgresHotelDAO) GetByAttribute(attrName string, attrValue string) (resLst list.List, err error) {
	db, err := sql.Open(`postgres`, dao.connStr)

	if err != nil {
		return resLst, errors.New(serverrors.ErrDatabaseConnection)
	}

	rows, err := db.Query(
		`select * from hotels where $1 = '$2';`,
		attrName, attrValue,
	)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return resLst, errors.New(serverrors.ErrQueryResRead)
		}

		return resLst, nil
	}

	for rows.Next() {
		var hotel models.Hotel
		err = rows.Scan(&hotel)

		if err != nil {
			return list.List{}, errors.New(serverrors.ErrQueryResRead)
		}

		resLst.PushBack(hotel)
	}

	return resLst, nil
}

func (dao *PostgresHotelDAO) Update(hotel *models.Hotel) (models.Hotel, error) {
	return models.Hotel{}, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresHotelDAO) Delete(hotel *models.Hotel) error {
	return errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresHotelDAO) DeleteByAttr(attrName string, attrValue string) error {
	return errors.New(serverrors.ErrMethodIsNotImplemented)
}
