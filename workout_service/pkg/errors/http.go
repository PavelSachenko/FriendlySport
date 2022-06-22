package errors

import "errors"

var (
	NotFound            = errors.New("Not Found")
	BadRequest          = errors.New("Bad Request")
	Unauthorized        = errors.New("Unauthorized")
	UnprocessableEntity = errors.New("Unprocessable Entity")
)
