package handler

import (
	"backend_test/model"
	"backend_test/service"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	employeeService service.EmployeeService
}

func NewHandler(
	employeeService service.EmployeeService,
) *Handler {
	return &Handler{
		employeeService: employeeService,
	}
}

func defaultPageRequest(pr *model.PageRequest) {
	if pr.PageNum == 0 {
		pr.PageNum = 1
	}
	if pr.PageSize == 0 {
		pr.PageSize = 10
	}
}

func RegisterHandlers(e *echo.Echo, h *Handler) {

	e.GET("/employees", h.GetEmployees)
	e.GET("/employees/:id", h.GetEmployeeByID)
	e.POST("/employees", h.AddEmployee)
	e.PUT("/employees/:id", h.EditEmployee)
	e.DELETE("/employees/:id", h.DeleteEmployeeByID)

}
