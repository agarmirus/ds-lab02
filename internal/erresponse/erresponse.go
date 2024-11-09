package erresponse

import (
	"fmt"
	"strings"
)

type ErrorResponse struct {
	msg  string
	errs map[string]string
}

func (erres *ErrorResponse) SetMessage(msg string) {
	erres.msg = msg
}

func (erres *ErrorResponse) GetMessage() string {
	return erres.msg
}

func (erres *ErrorResponse) AddFieldError(field string, err string) {
	erres.errs[field] = err
}

func (erres *ErrorResponse) ClearFieldErrors() {
	erres.errs = make(map[string]string)
}

func (erres *ErrorResponse) ToString() string {
	var result string = ""

	if strings.TrimSpace(erres.msg) != "" {
		result = fmt.Sprintf(`{"message":"%s"`, erres.msg)

		if len(erres.errs) > 0 {
			errorStrs := make([]string, 0, len(erres.errs))

			for field, err := range erres.errs {
				errorStrs = append(
					errorStrs,
					fmt.Sprintf(`{"field":"%s","error":"%s"}`, field, err),
				)
			}

			result += `,"errors: ["` + strings.Join(errorStrs, ",") + "]"
		}

		result += "}"
	}

	return result
}
