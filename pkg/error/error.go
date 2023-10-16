package pkgerror

import (
	"net/http"
)

type CustomError struct {
	Msg      string
	Code     string
	HttpCode int
	Err      error
}

func (e CustomError) WithError(err error) CustomError {
	if err != nil {
		e.Err = err
	}
	return e
}

func (e CustomError) IsNoError() bool {
	return e == CustomError{}
}

var (
	NoError                    CustomError = CustomError{}
	ErrSystemError             CustomError = CustomError{Code: "9999", Msg: "Unexpected error occured, please try again later", HttpCode: http.StatusInternalServerError}
	ErrUnauthorizedRequest     CustomError = CustomError{Code: "0001", Msg: "Request unauthorized. Missing or invalid access token", HttpCode: http.StatusUnauthorized}
	ErrInvalidParams           CustomError = CustomError{Code: "0002", Msg: "Missing or invalid request params, headers, or body", HttpCode: http.StatusBadRequest}
	ErrUndefinedPathPermission CustomError = CustomError{Code: "0003", Msg: "Undefined request path and/or permission mapping", HttpCode: http.StatusNotFound}
	ErrForbiddenRequest        CustomError = CustomError{Code: "0004", Msg: "Request forbidden. Operation not allowed", HttpCode: http.StatusForbidden}
	ErrEmployeeNotFound        CustomError = CustomError{Code: "0005", Msg: "Employee not found", HttpCode: http.StatusNotFound}
	ErrEmployeeIsExist         CustomError = CustomError{Code: "0006", Msg: "Employee is already exist", HttpCode: http.StatusBadRequest}
)
