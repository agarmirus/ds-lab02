package serverrors

const ErrMethodIsNotImplemented string = `method is not implemented`

const ErrDatabaseConnection string = `cannot connect to database`
const ErrEntityInsert string = `error while inserting new entity`
const ErrQueryResRead string = `error while reading SQL-query result`
const ErrQueryExec string = `error while executing SQL-query`

const ErrEntityNotFound string = `entity not found in database`

const ErrInvalidReservUid string = `invalid reservation UID field`
const ErrInvalidReservUsername string = `invalid reservation username field`
const ErrInvalidReservPayUID string = `invalid reservation payment UID field`
const ErrInvalidReservHotelId string = `invalid reservation hotel ID field`
const ErrInvalidReservStatus string = `invalid reservation status field`
const ErrInvalidReservDates string = `invalid reservation dates`

const ErrInvalidHotelId string = `invalid hotel ID field`

const ErrUnknown string = `unknown error`
