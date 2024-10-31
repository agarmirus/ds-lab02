package models

type Loyalty struct {
	id               int
	username         string
	reservationCount int
	status           string
	discount         int
}

func (loyalty *Loyalty) SetId(id int) {
	loyalty.id = id
}

func (loyalty *Loyalty) SetUsername(username string) {
	loyalty.username = username
}

func (loyalty *Loyalty) SetReservationCount(reservationCount int) {
	loyalty.reservationCount = reservationCount
}

func (loyalty *Loyalty) SetStatus(status string) {
	loyalty.status = status
}

func (loyalty *Loyalty) SetDiscount(discount int) {
	loyalty.discount = discount
}

func (loyalty *Loyalty) GetId() int {
	return loyalty.id
}

func (loyalty *Loyalty) GetUsername() string {
	return loyalty.username
}

func (loyalty *Loyalty) GetReservationCount() int {
	return loyalty.reservationCount
}

func (loyalty *Loyalty) GetStatus() string {
	return loyalty.status
}

func (loyalty *Loyalty) GetDiscount() int {
	return loyalty.discount
}
