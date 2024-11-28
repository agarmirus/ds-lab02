package serverrors

// Startup errors
const ErrConfigRead string = `error while reading config file`
const ErrServiceBuild string = `error while building service`
const ErrControllerPrepare string = `error while preparing controller`

// Implementation errors
const ErrMethodIsNotImplemented string = `method is not implemented`

// Database errors
const ErrDatabaseConnection string = `cannot connect to database`
const ErrEntityInsert string = `error while inserting new entity`
const ErrQueryResRead string = `error while reading SQL-query result`
const ErrQueryExec string = `error while executing SQL-query`

// Data errors
const ErrInvalidPagesData string = `invalid pages data`
const ErrInvalidUsername string = `invalid username`
const ErrInvalidCrReservReq string = `invalid create reservation request data`

const ErrInvalidReservUid string = `invalid reservation UID`

const ErrInvalidReservUsername string = `invalid reservation username field`
const ErrInvalidReservPayUID string = `invalid reservation payment UID field`
const ErrInvalidReservHotelId string = `invalid reservation hotel ID field`
const ErrInvalidReservStatus string = `invalid reservation status field`
const ErrInvalidReservDates string = `invalid reservation dates`

const ErrInvalidHotelId string = `invalid hotel ID field`
const ErrInvalidHoteUid string = `invalid hotel UID`

// Result errors
const ErrEntityNotFound string = `entity not found in database`
const ErrReservNotFound string = `reservation not found`
const ErrHotelNotFound string = `hotel not found`
const ErrPaymentNotFound string = `payment not found`
const ErrLoyaltyNotFound string = `loyalty not found`

// HTTP errors
const ErrNewRequestForming string = `error while creating new request`
const ErrRequestSend string = `error while sending request to service`
const ErrResponseRead string = `error while reading service response`
const ErrResponseParse string = `error while parsing service response`

// Internal errors
const ErrJSONParse string = `error while writting entity into json`

// Unknown :P
const ErrUnknown string = `unknown error`
