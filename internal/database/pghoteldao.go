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
	db, localErr := sql.Open(`postgres`, dao.connStr)

	if localErr == nil {
		var rows *sql.Rows
		rows, localErr = db.Query(`select * from hotels;`)
		for localErr == nil && rows.Next() {
			var hotel models.Hotel
			localErr = rows.Scan(&hotel)

			if localErr != nil {
				err = errors.New(serverrors.ErrQueryResRead)
			}

			resLst.PushBack(hotel)
		}
	} else {
		err = errors.New(serverrors.ErrDatabaseConnection)
	}

	return resLst, err
}

func (dao *PostgresHotelDAO) GetById(hotel *models.Hotel) (resHotel models.Hotel, err error) {
	if hotel.Id <= 0 {
		err = errors.New(serverrors.ErrInvalidHotelId)
	} else {
		db, localErr := sql.Open(`postgres`, dao.connStr)

		if localErr == nil {
			row := db.QueryRow(
				`select * from hotels where id = $1;`,
				hotel.Id,
			)
			localErr = row.Scan(&resHotel)

			if localErr != nil {
				if errors.Is(localErr, sql.ErrNoRows) {
					err = errors.New(serverrors.ErrEntityNotFound)
				} else {
					err = errors.New(serverrors.ErrQueryResRead)
				}
			}
		} else {
			err = errors.New(serverrors.ErrDatabaseConnection)
		}
	}

	return resHotel, err
}

func (dao *PostgresHotelDAO) GetByAttribute(attrName string, attrValue string) (resLst list.List, err error) {
	db, localErr := sql.Open(`postgres`, dao.connStr)

	if localErr == nil {
		var rows *sql.Rows
		rows, localErr = db.Query(
			`select * from hotels where $1 = '$2';`,
			attrName, attrValue,
		)
		for localErr == nil && rows.Next() {
			var hotel models.Hotel
			localErr = rows.Scan(&hotel)

			if localErr != nil {
				err = errors.New(serverrors.ErrQueryResRead)
			}

			resLst.PushBack(hotel)
		}
	} else {
		err = errors.New(serverrors.ErrDatabaseConnection)
	}

	return resLst, err
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
