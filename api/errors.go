package api

import "net/http"

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

//Error impolements the error interface
func (e Error) Error() string{
	return e.Err
}

func NewError(code int, error string) Error {
	return Error{
		Code: code,
		Err:  error,
	}
}

func ErrUnauthorized() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err: "unauthorized request",
	}
}

func ErrInvalidID() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err: "invalid id is given",
	}
}