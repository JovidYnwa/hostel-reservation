package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiError, ok := err.(Error); ok {
		return c.Status(apiError.Code).JSON(apiError)
	}
	apiError := NewError(http.StatusInternalServerError, err.Error())
	return c.Status(apiError.Code).JSON(apiError)
}

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

func ErrBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err: "Ivalid JSON request",
	}
}

func ErrNotResourceNotFound(res string) Error {
	return Error{
		Code: http.StatusNotFound,
		Err: res + "not found",
	}
}

func ErrInvalidID() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err: "invalid id is given",
	}
}