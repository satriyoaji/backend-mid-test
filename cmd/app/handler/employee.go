package handler

import (
	"backend_test/model"
	"backend_test/pkg/util/responseutil"
	"backend_test/pkg/validator"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetEmployees(ctx echo.Context) error {
	var filter model.GetEmployeesFilter
	if err := validator.BindAndValidate(ctx, &filter); !err.IsNoError() {
		return responseutil.SendErrorResponse(ctx, err)
	}
	defaultPageRequest(&filter.PageRequest)
	results, err := h.employeeService.GetEmployees(ctx, filter)
	if err.IsNoError() {
		return responseutil.SendSuccessReponse(ctx, results, nil)
	}
	return responseutil.SendErrorResponse(ctx, err)
}

func (h *Handler) AddEmployee(ctx echo.Context) error {
	req := model.CreateEmployeeRequest{}
	if err := validator.BindAndValidate(ctx, &req); !err.IsNoError() {
		return responseutil.SendErrorResponse(ctx, err)
	}
	result, ce := h.employeeService.CreateEmployee(ctx, req)
	if ce.IsNoError() {
		return responseutil.SendSuccessReponse(ctx, result, nil)
	}
	return responseutil.SendErrorResponse(ctx, ce)
}

func (h *Handler) GetEmployeeByID(ctx echo.Context) error {
	req := model.GetEmployeeByIDRequest{}
	if err := validator.BindAndValidate(ctx, &req); !err.IsNoError() {
		return responseutil.SendErrorResponse(ctx, err)
	}
	result, ce := h.employeeService.GetEmployeeByID(ctx, req)
	if ce.IsNoError() {
		return responseutil.SendSuccessReponse(ctx, result, nil)
	}
	return responseutil.SendErrorResponse(ctx, ce)
}

func (h *Handler) EditEmployee(ctx echo.Context) error {
	req := model.EditEmployeeRequest{}
	if err := validator.BindAndValidate(ctx, &req); !err.IsNoError() {
		return responseutil.SendErrorResponse(ctx, err)
	}
	result, ce := h.employeeService.EditEmployee(ctx, req)
	if ce.IsNoError() {
		return responseutil.SendSuccessReponse(ctx, result, nil)
	}
	return responseutil.SendErrorResponse(ctx, ce)
}

func (h *Handler) DeleteEmployeeByID(ctx echo.Context) error {
	req := model.DeleteEmployeeByIDRequest{}
	if err := validator.BindAndValidate(ctx, &req); !err.IsNoError() {
		return responseutil.SendErrorResponse(ctx, err)
	}
	ce := h.employeeService.DeleteEmployeeByID(ctx, req)
	if ce.IsNoError() {
		return responseutil.SendSuccessReponse(ctx, nil, nil)
	}
	return responseutil.SendErrorResponse(ctx, ce)
}
