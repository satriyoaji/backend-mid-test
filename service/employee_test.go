package service

import (
	"backend_test/entity"
	mocks "backend_test/mocks/repository"
	"backend_test/model"
	pkgerror "backend_test/pkg/error"
	"backend_test/pkg/util/copyutil"
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
)

func getExpectedEmployeesResult() *[]model.GetEmployeesResult {
	productResult := model.GetEmployeesResult{
		ID:        1,
		FirstName: "First Employee 0",
		LastName:  "Last Name 0",
		Email:     "employee@email.com",
	}
	return &[]model.GetEmployeesResult{productResult}
}

func TestGetEmployees(t *testing.T) {
	testCases := []struct {
		Name           string
		InitService    func(r *mocks.Repository) EmployeeService
		Context        echo.Context
		RequestParam   model.GetEmployeesFilter
		ExpectedCount  int
		ExpectedResult *[]model.GetEmployeesResult
		ExpectedError  pkgerror.CustomError
	}{
		{
			Name: "FindAllEmployeesError",
			InitService: func(r *mocks.Repository) EmployeeService {
				r.On("FindAllEmployees", context.Background(), model.GetEmployeesFilter{}).Return([]entity.Employee{}, errors.New("database error"))
				return NewEmployeeService(r)
			},
			Context:        createEchoContext(false),
			RequestParam:   model.GetEmployeesFilter{},
			ExpectedCount:  len(*getExpectedEmployeesResult()),
			ExpectedResult: nil,
			ExpectedError:  pkgerror.ErrSystemError,
		},
		{
			Name: "FindAllEmployeesSuccess",
			InitService: func(r *mocks.Repository) EmployeeService {
				employee := entity.Employee{
					ID:        1,
					FirstName: "First Employee 0",
					LastName:  "Last Name 0",
					Email:     "employee@email.com",
				}
				expectedReturn := []entity.Employee{employee}
				r.On("FindAllEmployees", context.Background(), model.GetEmployeesFilter{}).Return(expectedReturn, nil)
				return NewEmployeeService(r)
			},
			Context:        createEchoContext(false),
			RequestParam:   model.GetEmployeesFilter{},
			ExpectedCount:  len(*getExpectedEmployeesResult()),
			ExpectedResult: getExpectedEmployeesResult(),
			ExpectedError:  pkgerror.NoError,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			r := new(mocks.Repository)
			s := tc.InitService(r)
			results, err := s.GetEmployees(tc.Context, tc.RequestParam)
			assert.Equal(t, tc.ExpectedResult, results)
			assert.Equal(t, tc.ExpectedError.Code, err.Code)
			assert.Equal(t, tc.ExpectedError.HttpCode, err.HttpCode)
			assert.Equal(t, tc.ExpectedError.Msg, err.Msg)
			if tc.ExpectedResult != nil {
				assert.NotNil(t, results)
				assert.Equal(t, tc.ExpectedResult, results)
			}
			r.AssertExpectations(t)
		})
	}
}

func TestGetEmployeeByID(t *testing.T) {
	employee := entity.Employee{
		ID: 1,
	}
	result := model.GetEmployeeByIDResult{}
	copyutil.Copy(&employee, &result)
	testCases := []struct {
		Name           string
		InitService    func(r *mocks.Repository) EmployeeService
		Context        echo.Context
		ExpectedResult *model.GetEmployeeByIDResult
		ExpectedError  pkgerror.CustomError
	}{
		{
			Name: "SystemError",
			InitService: func(r *mocks.Repository) EmployeeService {
				r.On("FindEmployeeByID", mock.Anything, mock.Anything).Return(entity.Employee{}, errors.New("database error"))
				return NewEmployeeService(r)
			},
			Context:       createEchoContext(true),
			ExpectedError: pkgerror.ErrSystemError,
		},
		{
			Name: "EmployeeNotFound",
			InitService: func(r *mocks.Repository) EmployeeService {
				r.On("FindEmployeeByID", mock.Anything, mock.Anything).Return(entity.Employee{}, gorm.ErrRecordNotFound)
				return NewEmployeeService(r)
			},
			Context:       createEchoContext(true),
			ExpectedError: pkgerror.ErrEmployeeNotFound,
		},
		{
			Name: "Success",
			InitService: func(r *mocks.Repository) EmployeeService {
				r.On("FindEmployeeByID", mock.Anything, mock.Anything).Return(employee, nil)
				return NewEmployeeService(r)
			},
			Context:        createEchoContext(true),
			ExpectedError:  pkgerror.NoError,
			ExpectedResult: &result,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			r := new(mocks.Repository)
			s := tc.InitService(r)
			result, err := s.GetEmployeeByID(tc.Context, model.GetEmployeeByIDRequest{})
			assert.Equal(t, tc.ExpectedError.Code, err.Code)
			assert.Equal(t, tc.ExpectedError.HttpCode, err.HttpCode)
			assert.Equal(t, tc.ExpectedError.Msg, err.Msg)
			if tc.ExpectedResult != nil {
				assert.NotNil(t, result)
				assert.Equal(t, tc.ExpectedResult, result)
			}
			r.AssertExpectations(t)
		})
	}
}

func TestCreateEmployee(t *testing.T) {
	testCases := []struct {
		Name           string
		InitService    func(r *mocks.Repository) EmployeeService
		Context        echo.Context
		Request        model.CreateEmployeeRequest
		ExpectedResult *model.CreateEmployeeResult
		ExpectedError  pkgerror.CustomError
	}{
		{
			Name: "FindEmployeeByEmailErrorSystem",
			InitService: func(r *mocks.Repository) EmployeeService {
				r.On("FindEmployeeByEmail", context.Background(), "employee@email.com").Return(entity.Employee{}, errors.New("database error"))
				return NewEmployeeService(r)
			},
			Context: createEchoContext(false),
			Request: model.CreateEmployeeRequest{
				Email: "employee@email.com",
			},
			ExpectedResult: nil,
			ExpectedError:  pkgerror.ErrSystemError,
		},
		{
			Name: "FindEmployeeByEmailErrorExisted",
			InitService: func(r *mocks.Repository) EmployeeService {
				r.On("FindEmployeeByEmail", context.Background(), "employee@email.com").Return(entity.Employee{ID: uint(1), Email: "employee@email.com"}, nil)
				return NewEmployeeService(r)
			},
			Context: createEchoContext(false),
			Request: model.CreateEmployeeRequest{
				Email: "employee@email.com",
			},
			ExpectedResult: nil,
			ExpectedError:  pkgerror.ErrEmployeeIsExist,
		},
		{
			Name: "CreateEmployeeError",
			InitService: func(r *mocks.Repository) EmployeeService {
				r.On("FindEmployeeByEmail", context.Background(), "employee@email.com").Return(entity.Employee{}, nil)
				r.On("TxBegin").Return(nil)
				r.On("CreateEmployee", context.Background(), mock.Anything).Return(errors.New("database error"))
				r.On("TxRollback").Return(nil)
				return NewEmployeeService(r)
			},
			Context: createEchoContext(true),
			Request: model.CreateEmployeeRequest{
				FirstName: "First Employee 0",
				LastName:  "Last Name 0",
				Email:     "employee@email.com",
				HireDate:  "2023-09-20",
			},
			ExpectedResult: nil,
			ExpectedError:  pkgerror.ErrSystemError,
		},
		{
			Name: "CreateEmployeeSuccess",
			InitService: func(r *mocks.Repository) EmployeeService {
				r.On("FindEmployeeByEmail", context.Background(), "employee@email.com").Return(entity.Employee{}, nil)
				r.On("TxBegin").Return(nil)
				r.On("CreateEmployee", context.Background(), mock.Anything).Return(nil)
				r.On("TxCommit").Return(nil)
				return NewEmployeeService(r)
			},
			Context: createEchoContext(true),
			Request: model.CreateEmployeeRequest{
				FirstName: "First Employee 0",
				LastName:  "Last Name 0",
				Email:     "employee@email.com",
			},
			ExpectedResult: &model.CreateEmployeeResult{
				FirstName: "First Employee 0",
				LastName:  "Last Name 0",
				Email:     "employee@email.com",
			},
			ExpectedError: pkgerror.NoError,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			r := new(mocks.Repository)
			s := tc.InitService(r)
			result, err := s.CreateEmployee(tc.Context, tc.Request)
			assert.Equal(t, tc.ExpectedError.Code, err.Code)
			assert.Equal(t, tc.ExpectedError.HttpCode, err.HttpCode)
			assert.Equal(t, tc.ExpectedError.Msg, err.Msg)
			if tc.ExpectedResult != nil {
				assert.NotNil(t, result)
				assert.Equal(t, tc.ExpectedResult, result)
			}
			r.AssertExpectations(t)
		})
	}
}
