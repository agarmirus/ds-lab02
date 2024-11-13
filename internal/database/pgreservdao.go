package database

import (
	"container/list"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx"

	"github.com/agarmirus/ds-lab02/internal/models"
	"github.com/agarmirus/ds-lab02/internal/serverrors"
)

type PostgresReservationDAO struct {
	connStr string
}

func (dao *PostgresReservationDAO) SetConnectionString(connStr string) {
	dao.connStr = connStr
}

func validateReservation(reservation *models.Reservation) (err error) {
	if uuid.Validate(reservation.Uid) != nil {
		err = errors.New(serverrors.ErrInvalidReservUid)
	} else if strings.Trim(reservation.Username, ` `) == `` {
		err = errors.New(serverrors.ErrInvalidReservUsername)
	} else if uuid.Validate(reservation.Uid) != nil {
		err = errors.New(serverrors.ErrInvalidReservPayUID)
	} else if reservation.HotelId <= 0 {
		err = errors.New(serverrors.ErrInvalidReservHotelId)
	} else if reservation.Status != `PAID` && reservation.Status != `CANCELED` {
		err = errors.New(serverrors.ErrInvalidReservStatus)
	} else if reservation.StartDate.Unix() < reservation.EndDate.Unix() {
		err = errors.New(serverrors.ErrInvalidReservDates)
	}

	return err
}

func (dao *PostgresReservationDAO) Create(reservation *models.Reservation) (newReservation models.Reservation, err error) {
	err = validateReservation(reservation)

	if err == nil {
		db, localErr := sql.Open(`postgres`, dao.connStr)

		if localErr == nil {
			row := db.QueryRow(
				`insert into reservation (reservation_uid, username, payment_uid, hotel_id, status, start_date, end_date)
				values ('$1', '$2', '$3', $4, '$5', '$6', '$7')
				returning *`,
				reservation.Uid, reservation.Username,
				reservation.PaymentUid, reservation.HotelId,
				reservation.Status, reservation.StartDate.Format(`%F`),
				reservation.EndDate.Format(`%F`),
			)
			localErr = row.Scan(&newReservation)

			if localErr != nil {
				if errors.Is(localErr, sql.ErrNoRows) {
					err = errors.New(serverrors.ErrEntityInsert)
				} else {
					err = errors.New(serverrors.ErrQueryResRead)
				}
			}
		} else {
			err = errors.New(serverrors.ErrDatabaseConnection)
		}
	}

	return newReservation, err
}

func (dao *PostgresReservationDAO) Get() (list.List, error) {
	return list.List{}, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresReservationDAO) GetById(reservation *models.Reservation) (models.Reservation, error) {
	return models.Reservation{}, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresReservationDAO) GetByAttribute(attrName string, attrValue string) (resLst list.List, err error) {
	db, localErr := sql.Open(`postgres`, dao.connStr)

	if localErr == nil {
		var rows *sql.Rows
		rows, localErr = db.Query(
			`select * from reservation where $1 == '$2'`,
			attrName, attrValue,
		)
		for localErr == nil && rows.Next() {
			var reservation models.Reservation
			localErr = rows.Scan(&reservation)

			if localErr != nil {
				err = errors.New(serverrors.ErrQueryResRead)
			}

			resLst.PushBack(reservation)
		}
	} else {
		err = errors.New(serverrors.ErrDatabaseConnection)
	}

	return resLst, err
}

func (dao *PostgresReservationDAO) Update(reservation *models.Reservation) (models.Reservation, error) {
	return models.Reservation{}, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresReservationDAO) Delete(reservation *models.Reservation) error {
	return errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresReservationDAO) DeleteByAttr(attrName string, attrValue string) (err error) {
	db, localErr := sql.Open(`postgres`, dao.connStr)

	if localErr == nil {
		_, localErr = db.Exec(
			`delete from reservation where $1 = '$2'`,
			attrName, attrValue,
		)

		if localErr != nil {
			err = errors.New(serverrors.ErrQueryExec)
		}
	}

	return err
}
