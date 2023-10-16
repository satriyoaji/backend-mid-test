package service

import (
	"backend_test/entity"
	"backend_test/model"
	"backend_test/repository"
	"errors"

	pkgerror "backend_test/pkg/error"
	"backend_test/pkg/util/copyutil"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type EmployeeService interface {
	GetEmployees(ctx echo.Context, filter model.GetEmployeesFilter) (*[]model.GetEmployeesResult, pkgerror.CustomError)
	CreateEmployee(ctx echo.Context, req model.CreateEmployeeRequest) (*model.CreateEmployeeResult, pkgerror.CustomError)
	GetEmployeeByID(ctx echo.Context, req model.GetEmployeeByIDRequest) (*model.GetEmployeeByIDResult, pkgerror.CustomError)
	EditEmployee(ctx echo.Context, req model.EditEmployeeRequest) (*model.EditEmployeeResult, pkgerror.CustomError)
	DeleteEmployeeByID(ctx echo.Context, req model.DeleteEmployeeByIDRequest) pkgerror.CustomError
}

type EmployeeServiceImpl struct {
	repo repository.Repository
}

func NewEmployeeService(
	repo repository.Repository) *EmployeeServiceImpl {
	return &EmployeeServiceImpl{
		repo: repo,
	}
}

func (s EmployeeServiceImpl) GetEmployees(ctx echo.Context, filter model.GetEmployeesFilter) (*[]model.GetEmployeesResult, pkgerror.CustomError) {
	rctx := ctx.Request().Context()
	results := []model.GetEmployeesResult{}
	employees, err := s.repo.FindAllEmployees(rctx, filter)
	if err != nil {
		log.Error("Find employees error: ", err)
		return nil, pkgerror.ErrSystemError
	}
	copyutil.Copy(&employees, &results)
	return &results, pkgerror.NoError
}

func (s *EmployeeServiceImpl) GetEmployeeByID(ctx echo.Context, req model.GetEmployeeByIDRequest) (*model.GetEmployeeByIDResult, pkgerror.CustomError) {
	rctx := ctx.Request().Context()
	employee, err := s.repo.FindEmployeeByID(rctx, uint(req.EmployeeID))
	if err != nil {
		log.Error("Find employee by ID error: ", err)
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, pkgerror.ErrEmployeeNotFound.WithError(err)
		}
		return nil, pkgerror.ErrSystemError.WithError(err)
	}
	result := model.GetEmployeeByIDResult{}
	copyutil.Copy(&employee, &result)
	return &result, pkgerror.NoError
}

func (s *EmployeeServiceImpl) CreateEmployee(ctx echo.Context, req model.CreateEmployeeRequest) (*model.CreateEmployeeResult, pkgerror.CustomError) {
	rctx := ctx.Request().Context()

	employeeFound, err := s.repo.FindEmployeeByEmail(rctx, req.Email)
	if err != nil {
		log.Error("Find user by Email error: ", err)
		if !errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, pkgerror.ErrSystemError.WithError(err)
		}
	}
	if employeeFound.Email != "" {
		return nil, pkgerror.ErrEmployeeIsExist.WithError(errors.New("Employee `email` is already created."))
	}

	txSuccess := false
	err = s.repo.TxBegin()
	if err != nil {
		log.Error("Start db transaction error: ", err)
		return nil, pkgerror.ErrSystemError.WithError(err)
	}
	defer func() {
		if r := recover(); r != nil || !txSuccess {
			err = s.repo.TxRollback()
			if err != nil {
				log.Error("Rollback db transaction error: ", err)
			}
		}
	}()

	var employee entity.Employee
	copyutil.Copy(&req, &employee)
	err = s.repo.CreateEmployee(rctx, &employee)
	if err != nil {
		log.Error("Create employee error: ", err)
		return nil, pkgerror.ErrSystemError.WithError(err)
	}
	err = s.repo.TxCommit()
	if err != nil {
		log.Error("Commit db transaction error: ", err)
	}
	var result model.CreateEmployeeResult
	copyutil.Copy(&employee, &result)
	txSuccess = true
	return &result, pkgerror.NoError
}

func (s *EmployeeServiceImpl) EditEmployee(ctx echo.Context, req model.EditEmployeeRequest) (*model.EditEmployeeResult, pkgerror.CustomError) {
	rctx := ctx.Request().Context()
	employee, err := s.repo.FindEmployeeByID(rctx, uint(req.EmployeeID))
	if err != nil {
		log.Error("Find employee by ID error: ", err)
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, pkgerror.ErrEmployeeNotFound.WithError(err)
		}
		return nil, pkgerror.ErrSystemError.WithError(err)
	}
	// Start transaction
	txSuccess := false
	err = s.repo.TxBegin()
	if err != nil {
		log.Error("Start db transaction error: ", err)
		return nil, pkgerror.ErrSystemError.WithError(err)
	}
	defer func() {
		if r := recover(); r != nil || !txSuccess {
			err = s.repo.TxRollback()
			if err != nil {
				log.Error("Rollback db transaction error: ", err)
			}
		}
	}()
	err = s.repo.UpdateEmployee(rctx, &employee)
	if err != nil {
		return nil, pkgerror.ErrSystemError.WithError(err)
	}
	// Commit transaction
	err = s.repo.TxCommit()
	if err != nil {
		log.Error("Commit db transaction error: ", err)
	}
	result := model.EditEmployeeResult{}
	copyutil.Copy(&employee, &result)
	txSuccess = true
	return &result, pkgerror.NoError
}

func (s *EmployeeServiceImpl) DeleteEmployeeByID(ctx echo.Context, req model.DeleteEmployeeByIDRequest) pkgerror.CustomError {
	rctx := ctx.Request().Context()
	err := s.repo.DeleteEmployee(rctx, uint(req.EmployeeID))
	if err != nil {
		log.Error("Delete employee by ID error: ", err)
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return pkgerror.ErrEmployeeNotFound.WithError(err)
		}
		return pkgerror.ErrSystemError.WithError(err)
	}

	return pkgerror.NoError
}
