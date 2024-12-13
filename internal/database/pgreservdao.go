package database

import (
	"container/list"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

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
		err = serverrors.ErrInvalidReservUid
	} else if strings.Trim(reservation.Username, ` `) == `` {
		err = serverrors.ErrInvalidReservUsername
	} else if uuid.Validate(reservation.Uid) != nil {
		err = serverrors.ErrInvalidReservPayUID
	} else if reservation.HotelId <= 0 {
		err = serverrors.ErrInvalidReservHotelId
	} else if reservation.Status != `PAID` && reservation.Status != `CANCELED` {
		err = serverrors.ErrInvalidReservStatus
	} else if reservation.StartDate > reservation.EndDate {
		err = serverrors.ErrInvalidReservDates
	}

	return err
}

func (dao *PostgresReservationDAO) Create(reservation *models.Reservation) (newReservation models.Reservation, err error) {
	err = validateReservation(reservation)

	if err != nil {
		log.Println("[ERROR] PostgresReservationDAO.Create. Invalid reservation data:", err)
		return newReservation, err
	}

	db, err := sql.Open(`postgres`, dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresReservationDAO.Create. Cannot connect to database:", err)
		return newReservation, serverrors.ErrDatabaseConnection
	}

	defer db.Close()

	row := db.QueryRow(
		`insert into reservation (reservation_uid, username, payment_uid, hotel_id, status, start_date, end_date)
		values ($1, $2, $3, $4, $5, $6, $7)
		returning *;`,
		reservation.Uid, reservation.Username,
		reservation.PaymentUid, reservation.HotelId,
		reservation.Status, reservation.StartDate,
		reservation.EndDate,
	)
	err = row.Scan(
		&newReservation.Id, &newReservation.Uid,
		&newReservation.Username, &newReservation.PaymentUid,
		&newReservation.HotelId, &newReservation.Status,
		&newReservation.StartDate, &newReservation.EndDate,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[ERROR] PostgresReservationDAO.Create. Entity not found")
			err = serverrors.ErrEntityNotFound
		} else {
			log.Println("[ERROR] PostgresReservationDAO.Create. Error while reading query result:", err)
			err = serverrors.ErrQueryResRead
		}
	}

	return newReservation, err
}

func (dao *PostgresReservationDAO) Get() (list.List, error) {
	log.Println("[ERROR] PostgresReservationDAO.Get. Method is not implemented")
	return list.List{}, serverrors.ErrMethodIsNotImplemented
}

func (dao *PostgresReservationDAO) GetPaginated(
	page int,
	pageSize int,
) (resLst list.List, err error) {
	log.Println("[ERROR] PostgresLoyaltyDAO.GetPaginated. Method is not implemented")
	return list.List{}, serverrors.ErrMethodIsNotImplemented
}

func (dao *PostgresReservationDAO) GetById(reservation *models.Reservation) (models.Reservation, error) {
	log.Println("[ERROR] PostgresReservationDAO.GetById. Method is not implemented")
	return models.Reservation{}, serverrors.ErrMethodIsNotImplemented
}

func (dao *PostgresReservationDAO) GetByAttribute(attrName string, attrValue string) (resLst list.List, err error) {
	db, err := sql.Open(`postgres`, dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresReservationDAO.GetByAttribute. Cannot connect to database:", err)
		return resLst, serverrors.ErrDatabaseConnection
	}

	defer db.Close()

	queryStr := fmt.Sprintf(
		`select * from reservation where %s = $1;`,
		attrName,
	)
	rows, err := db.Query(queryStr, attrValue)

	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Println("[ERROR] PostgresReservationDAO.GetByAttribute. Error while executing query:", err)
			return resLst, serverrors.ErrQueryResRead
		}

		return resLst, nil
	}

	for rows.Next() {
		var reservation models.Reservation
		err = rows.Scan(
			&reservation.Id, &reservation.Uid,
			&reservation.Username, &reservation.PaymentUid,
			&reservation.HotelId, &reservation.Status,
			&reservation.StartDate, &reservation.EndDate,
		)

		if err != nil {
			log.Println("[ERROR] PostgresReservationDAO.GetByAttribute. Error while reading query result:", err)
			return list.List{}, serverrors.ErrQueryResRead
		}

		resLst.PushBack(reservation)
	}

	return resLst, nil
}

func (dao *PostgresReservationDAO) Update(reservation *models.Reservation) (updatedReservation models.Reservation, err error) {
	db, err := sql.Open(`postgres`, dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresReservationDAO.Update. Cannot connect to database:", err)
		return updatedReservation, serverrors.ErrDatabaseConnection
	}

	defer db.Close()

	log.Println("[TRACE] PostgresReservationDAO.Update. Status =", reservation.Status)

	row := db.QueryRow(
		`update reservation
		set username = $1, payment_uid = $2, hotel_id = $3, status = $4, start_date = $5, end_date = $6
		where reservation_uid = $7
		returning *`,
		reservation.Username, reservation.PaymentUid, reservation.HotelId,
		reservation.Status, reservation.StartDate, reservation.EndDate,
		reservation.Uid,
	)
	err = row.Scan(
		&updatedReservation.Id, &updatedReservation.Uid,
		&updatedReservation.Username, &updatedReservation.PaymentUid,
		&updatedReservation.HotelId, &updatedReservation.Status,
		&updatedReservation.StartDate, &updatedReservation.EndDate,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("[ERROR] PostgresReservationDAO.Update. Entity not found")
			err = serverrors.ErrEntityNotFound
		} else {
			log.Println("[ERROR] PostgresReservationDAO.Update. Error while reading query result:", err)
			err = serverrors.ErrQueryResRead
		}
	}

	return updatedReservation, err
}

func (dao *PostgresReservationDAO) Delete(reservation *models.Reservation) error {
	log.Println("[ERROR] PostgresReservationDAO.Delete. Method is not implemented")
	return serverrors.ErrMethodIsNotImplemented
}

func (dao *PostgresReservationDAO) DeleteByAttr(attrName string, attrValue string) (err error) {
	db, err := sql.Open(`postgres`, dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresReservationDAO.DeleteByAttr. Cannot connect to database:", err)
		return serverrors.ErrDatabaseConnection
	}

	defer db.Close()

	_, err = db.Exec(
		`delete from reservation where $1 = $2;`,
		attrName, attrValue,
	)

	if err != nil {
		log.Println("[ERROR] PostgresReservationDAO.DeleteByAttr. Error while executing query:", err)
		return serverrors.ErrQueryExec
	}

	return nil
}
