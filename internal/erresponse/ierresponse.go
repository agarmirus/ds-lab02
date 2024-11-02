package erresponse

type IErrorResponse interface {
	SetMessage(string)
	GetMessage() string

	AddFieldError(field string, err string)

	ClearFieldErrors()
}
