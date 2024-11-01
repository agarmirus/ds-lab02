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

type PostgresHotelDAO struct {
	connStr string
}

func NewPostgresHotelDAO() IDAO[models.Hotel] {
	return &PostgresHotelDAO{}
}

func (dao *PostgresHotelDAO) SetConnectionString(connStr string) {
	dao.connStr = connStr
}

func (dao *PostgresHotelDAO) Create(*models.Hotel) error {
	return errors.New("PostgresHotelDAO.Create() is not implemented")
}

func (dao *PostgresHotelDAO) Get() (list.List, error) {
	var err error = nil
	resLst := list.List{}

	db, err := sql.Open("postgres", dao.connStr)

	if err == nil {
		defer db.Close()

		queryString := "select * from hotels;"

		var rows *sql.Rows
		rows, err = db.Query(queryString)

		if err == nil {
			defer rows.Close()

			for rows.Next() {
				var id, stars, price int
				var name, country, city, address string
				var uid uuid.UUID

				err = rows.Scan(
					&id, &uid,
					&name, &country,
					&city, &address,
					&stars, &price,
				)

				if err == nil {
					curHotel := models.Hotel{}
					curHotel.SetId(id)
					curHotel.SetUid(uid)
					curHotel.SetName(name)
					curHotel.SetCountry(country)
					curHotel.SetCity(city)
					curHotel.SetAddress(address)
					curHotel.SetStars(stars)
					curHotel.SetPrice(price)

					resLst.PushBack(curHotel)
				} else if errors.Is(err, sql.ErrNoRows) {
					err = nil
				}
			}
		}
	}

	return resLst, err
}

func (dao *PostgresHotelDAO) GetById(hotel *models.Hotel) (models.Hotel, error) {
	var err error = nil
	result := models.Hotel{}

	db, err := sql.Open("postgres", dao.connStr)

	if err == nil {
		defer db.Close()

		var id int = hotel.GetId()
		var queryString string = fmt.Sprintf("select * from hotels where id = %d;", id)

		var row *sql.Row = db.QueryRow(queryString)

		var stars, price int
		var uid uuid.UUID
		var name, country, city, address string

		err = row.Scan(
			&id, &uid,
			&name, &country,
			&city, &address,
			&stars, &price,
		)

		if err == nil {
			result.SetId(id)
			result.SetUid(uid)
			result.SetName(name)
			result.SetCountry(country)
			result.SetCity(city)
			result.SetAddress(address)
			result.SetStars(stars)
			result.SetPrice(price)
		} else if errors.Is(err, sql.ErrNoRows) {
			err = nil
		}
	}

	return result, err
}

func (dao *PostgresHotelDAO) GetByAttribute(attrName string, attrValue string) (list.List, error) {
	var err error = nil
	resLst := list.List{}

	db, err := sql.Open("postgres", dao.connStr)

	if err == nil {
		defer db.Close()

		queryString := fmt.Sprintf("select * from hotels where %s = '%s';", attrName, attrValue)

		var rows *sql.Rows
		rows, err = db.Query(queryString)

		if err == nil {
			defer rows.Close()

			for rows.Next() {
				var id, stars, price int
				var name, country, city, address string
				var uid uuid.UUID

				err = rows.Scan(
					&id, &uid,
					&name, &country,
					&city, &address,
					&stars, &price,
				)

				if err == nil {
					curHotel := models.Hotel{}
					curHotel.SetId(id)
					curHotel.SetUid(uid)
					curHotel.SetName(name)
					curHotel.SetCountry(country)
					curHotel.SetCity(city)
					curHotel.SetAddress(address)
					curHotel.SetStars(stars)
					curHotel.SetPrice(price)

					resLst.PushBack(curHotel)
				} else if errors.Is(err, sql.ErrNoRows) {
					err = nil
				}
			}
		}
	}

	return resLst, err
}

func (dao *PostgresHotelDAO) Update(*models.Hotel) (models.Hotel, error) {
	return models.Hotel{}, errors.New("PostgresHotelDAO.Update is not implemented")
}

func (dao *PostgresHotelDAO) Delete(*models.Hotel) error {
	return errors.New("PostgresHotelDAO.Delete is not implemented")
}
