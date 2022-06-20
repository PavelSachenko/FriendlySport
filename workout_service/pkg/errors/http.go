package errors

import "errors"

var (
	NotFound            = errors.New("Not Found")
	BadRequest          = errors.New("Bad Request")
	Unauthorized        = errors.New("Unauthorized")
	UnprocessableEntity = errors.New("Unprocessable Entity")
)

type AppHTTPError struct {
	Code      int
	HttpError error
}

func (a AppHTTPError) Error() string {
	return a.HttpError.Error()
}

func HTTPError(code int, err error) error {
	return &AppHTTPError{
		Code:      code,
		HttpError: err,
	}
}
