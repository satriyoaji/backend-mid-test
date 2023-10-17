package service

import (
	"backend_test/entity"
	mocks "backend_test/mocks/repository"
	"backend_test/model"
	pkgerror "backend_test/pkg/error"
	"backend_test/pkg/util/copyutil"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
)

func getExpectedEmployeesResult() *[]model.GetEmployeesResult {
	productResult := model.GetEmployeesResult{
		ID:        1,
		FirstName: fmt.Sprintf("First Employee %d", 1),
		LastName:  fmt.Sprintf("Last %d", 1),
		Email:     "",
	}
	return &[]model.GetEmployeesResult{productResult}
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
