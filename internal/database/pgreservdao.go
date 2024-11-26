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

func NewPostgresReservationDAO(connStr string) IDAO[models.Reservation] {
	return &PostgresReservationDAO{connStr}
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
		return newReservation, err
	}

	db, err := sql.Open(`postgres`, dao.connStr)

	if err != nil {
		return newReservation, errors.New(serverrors.ErrDatabaseConnection)
	}

	row := db.QueryRow(
		`insert into reservation (reservation_uid, username, payment_uid, hotel_id, status, start_date, end_date)
		values ('$1', '$2', '$3', $4, '$5', '$6', '$7')
		returning *;`,
		reservation.Uid, reservation.Username,
		reservation.PaymentUid, reservation.HotelId,
		reservation.Status, reservation.StartDate.Format(`%F`),
		reservation.EndDate.Format(`%F`),
	)
	err = row.Scan(&newReservation)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.New(serverrors.ErrEntityInsert)
		} else {
			err = errors.New(serverrors.ErrQueryResRead)
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
	db, err := sql.Open(`postgres`, dao.connStr)

	if err != nil {
		return resLst, errors.New(serverrors.ErrDatabaseConnection)
	}

	rows, err := db.Query(
		`select * from reservation where $1 == '$2';`,
		attrName, attrValue,
	)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return resLst, errors.New(serverrors.ErrQueryResRead)
		}

		return resLst, nil
	}

	for rows.Next() {
		var reservation models.Reservation
		err = rows.Scan(&reservation)

		if err != nil {
			return list.List{}, errors.New(serverrors.ErrQueryResRead)
		}

		resLst.PushBack(reservation)
	}

	return resLst, nil
}

func (dao *PostgresReservationDAO) Update(reservation *models.Reservation) (updatedReservation models.Reservation, err error) {
	db, err := sql.Open(`postgres`, dao.connStr)

	if err != nil {
		return updatedReservation, errors.New(serverrors.ErrDatabaseConnection)
	}

	row := db.QueryRow(
		`update reservation
		set username = '$1', payment_uid = '$2', hotel_id = $3, status = '$4', start_date = '$5', end_date = '$6'
		where reservation_uid = '$7'
		returning *`,
		reservation.Username, reservation.PaymentUid, reservation.HotelId,
		reservation.Status, reservation.StartDate.Format(`%F`), reservation.EndDate.Format(`%F`),
		reservation.Uid,
	)
	err = row.Scan(&updatedReservation)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.New(serverrors.ErrEntityNotFound)
		} else {
			err = errors.New(serverrors.ErrQueryResRead)
		}
	}

	return updatedReservation, err
}

func (dao *PostgresReservationDAO) Delete(reservation *models.Reservation) error {
	return errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresReservationDAO) DeleteByAttr(attrName string, attrValue string) (err error) {
	db, err := sql.Open(`postgres`, dao.connStr)

	if err != nil {
		return errors.New(serverrors.ErrDatabaseConnection)
	}

	_, err = db.Exec(
		`delete from reservation where $1 = '$2';`,
		attrName, attrValue,
	)

	if err != nil {
		return errors.New(serverrors.ErrQueryExec)
	}

	return nil
}
