package database

import (
	"container/list"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/agarmirus/ds-lab02/internal/models"
	"github.com/google/uuid"
)

type PostgresReservationDAO struct {
	connStr string
}

func NewPostgresReservationDAO() IDAO[models.Reservation] {
	return &PostgresReservationDAO{}
}

func (dao *PostgresReservationDAO) SetConnectionString(connStr string) {
	dao.connStr = connStr
}

func (dao *PostgresReservationDAO) Create(reservation *models.Reservation) error {
	var err error = nil

	db, err := sql.Open("postgres", dao.connStr)

	if err == nil {
		defer db.Close()

		queryString := fmt.Sprintf(
			"insert into table reservation (reservation_uid, username, payment_uid, hotel_id, status, start_date, end_data) values ('%s', '%s', '%s', %d, '%s', '%s', '%s');",
			reservation.GetUid().String(),
			reservation.GetUsername(),
			reservation.GetPaymentUid().String(),
			reservation.GetHotelId(),
			reservation.GetStatus(),
			reservation.GetStartDate().String(),
			reservation.GetEndDate().String(),
		)

		_, err = db.Exec(queryString)
	}

	return err
}

func (dao *PostgresReservationDAO) Get() (list.List, error) {
	return list.List{}, errors.New("PostgresReservationDAO.Get() is not implemented")
}

func (dao *PostgresReservationDAO) GetById(reservation *models.Reservation) (models.Reservation, error) {
	var err error = nil
	result := models.Reservation{}

	db, err := sql.Open("postgres", dao.connStr)

	if err == nil {
		defer db.Close()

		var id int = reservation.GetId()
		queryString := fmt.Sprintf("select * from reservation where id = %d", id)

		var row *sql.Row = db.QueryRow(queryString)

		var hotelId int
		var username, status string
		var uid, paymentUid uuid.UUID
		var startDate, endDate time.Time

		err = row.Scan(
			&id, &uid,
			&username, &paymentUid,
			&hotelId, &status,
			&startDate, &endDate,
		)

		if err == nil {
			result.SetId(id)
			result.SetUid(uid)
			result.SetUsername(username)
			result.SetPaymentUid(paymentUid)
			result.SetHotelId(hotelId)
			result.SetStatus(status)
			result.SetStartDate(startDate)
			result.SetEndDate(endDate)
		} else if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
	}

	return result, err
}

func (dao *PostgresReservationDAO) GetByAttribute(attrName string, attrValue string) (list.List, error) {
	var err error = nil
	resLst := list.List{}

	db, err := sql.Open("postgres", dao.connStr)

	if err == nil {
		defer db.Close()

		queryString := fmt.Sprintf("select * from reservation where %s = '%s';", attrName, attrValue)

		var rows *sql.Rows
		rows, err = db.Query(queryString)

		if err == nil {
			defer rows.Close()

			for rows.Next() {
				var id, hotelId int
				var username, status string
				var uid, paymentUid uuid.UUID
				var startDate, endDate time.Time

				err = rows.Scan(
					&id, &uid,
					&username, &paymentUid,
					&hotelId, &status,
					&startDate, &endDate,
				)

				if err == nil {
					curReservation := models.Reservation{}
					curReservation.SetId(id)
					curReservation.SetUid(uid)
					curReservation.SetUsername(username)
					curReservation.SetPaymentUid(paymentUid)
					curReservation.SetHotelId(hotelId)
					curReservation.SetStatus(status)
					curReservation.SetStartDate(startDate)
					curReservation.SetEndDate(endDate)

					resLst.PushBack(curReservation)
				} else if errors.Is(err, sql.ErrNoRows) {
					err = nil
				}
			}
		}
	}

	return resLst, err
}

func (dao *PostgresReservationDAO) Update(*models.Reservation) (models.Reservation, error) {
	return models.Reservation{}, errors.New("PostgresReservationDAO.Update() is not implemented")
}

func (dao *PostgresReservationDAO) Delete(reservation *models.Reservation) error {
	var err error = nil

	db, err := sql.Open("postgres", dao.connStr)

	if err == nil {
		defer db.Close()

		queryString := fmt.Sprintf(
			"delete from reservation where reservation_uid = '%s' and username = '%s';",
			reservation.GetUid().String(),
			reservation.GetUsername(),
		)

		_, err = db.Exec(queryString)
	}

	return err
}
