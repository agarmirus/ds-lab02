package models

import "github.com/google/uuid"

type Hotel struct {
	id      int
	uid     uuid.UUID
	name    string
	country string
	city    string
	address string
	stars   int
	price   int
}

func (hotel *Hotel) SetId(id int) {
	hotel.id = id
}

func (hotel *Hotel) SetUid(uid uuid.UUID) {
	hotel.uid = uid
}

func (hotel *Hotel) SetName(name string) {
	hotel.name = name
}

func (hotel *Hotel) SetCountry(country string) {
	hotel.country = country
}

func (hotel *Hotel) SetCity(city string) {
	hotel.city = city
}

func (hotel *Hotel) SetAddress(address string) {
	hotel.address = address
}

func (hotel *Hotel) SetStars(stars int) {
	hotel.stars = stars
}

func (hotel *Hotel) SetPrice(price int) {
	hotel.price = price
}

func (hotel *Hotel) GetId() int {
	return hotel.id
}

func (hotel *Hotel) GetUid() uuid.UUID {
	return hotel.uid
}

func (hotel *Hotel) GetName() string {
	return hotel.name
}

func (hotel *Hotel) GetCountry() string {
	return hotel.country
}

func (hotel *Hotel) GetCity() string {
	return hotel.city
}

func (hotel *Hotel) GetAddress() string {
	return hotel.address
}

func (hotel *Hotel) GetStars() int {
	return hotel.stars
}

func (hotel *Hotel) GetPrice() int {
	return hotel.price
}
