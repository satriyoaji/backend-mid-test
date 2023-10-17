package handler

import (
	mocks "backend_test/mocks/service"
	"backend_test/model"
	pkgerror "backend_test/pkg/error"
	"backend_test/pkg/util/jsonutil"
	"backend_test/pkg/util/responseutil"
	pkgvalidator "backend_test/pkg/validator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetEmployeeByID(t *testing.T) {
	result := model.GetEmployeeByIDResult{
		ID: 1,
	}
	testCases := []struct {
		Name                 string
		InitHandler          func(ctx echo.Context, s *mocks.EmployeeService) *Handler
		PathEmployeeID       string
		ExpectedHttpCode     int
		ExpectedResponseBody model.ResponseBody
	}{
		{
			Name: "InvalidParams",
			InitHandler: func(ctx echo.Context, s *mocks.EmployeeService) *Handler {
				return NewHandler(s)
			},
			ExpectedHttpCode:     http.StatusBadRequest,
			ExpectedResponseBody: responseutil.CreateErrorResponse(pkgerror.ErrInvalidParams),
		},
		{
			Name: "ServiceError",
			InitHandler: func(ctx echo.Context, s *mocks.EmployeeService) *Handler {
				s.On("GetEmployeeByID", mock.Anything, mock.Anything).Return(nil, pkgerror.ErrSystemError)
				return NewHandler(s)
			},
			PathEmployeeID:       "1",
			ExpectedHttpCode:     http.StatusInternalServerError,
			ExpectedResponseBody: responseutil.CreateErrorResponse(pkgerror.ErrSystemError),
		},
		{
			Name: "Success",
			InitHandler: func(ctx echo.Context, s *mocks.EmployeeService) *Handler {
				s.On("GetEmployeeByID", mock.Anything, mock.Anything).Return(&result, pkgerror.NoError)
				return NewHandler(s)
			},
			PathEmployeeID:       "1",
			ExpectedHttpCode:     http.StatusOK,
			ExpectedResponseBody: responseutil.CreateSuccessResponse(&result, nil),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			e := echo.New()
			e.Validator = pkgvalidator.New(validator.New())
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			c.SetPath("/v1/employees/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.PathEmployeeID)
			s := new(mocks.EmployeeService)
			h := tc.InitHandler(c, s)
			if assert.NoError(t, h.GetEmployeeByID(c)) {
				assert.Equal(t, tc.ExpectedHttpCode, res.Code)
				expected := tc.ExpectedResponseBody
				jsonpath, err := jsonutil.NewJsonPath(res.Body.String())
				assert.Nil(t, err)
				assert.Equal(t, expected.Status, jsonpath.GetString("status"))
				assert.Equal(t, expected.Code, jsonpath.GetString("code"))
				assert.Equal(t, expected.ErrorMessage, jsonpath.GetStringPtr("error_message"))
				if expected.Data != nil {
					data := expected.Data.(*model.GetEmployeeByIDResult)
					assert.Equal(t, data.ID, jsonpath.GetInt("data.id"))
				}
			}
			s.AssertExpectations(t)
		})
	}
}

func TestAddEmployee(t *testing.T) {
	validJson := `{
		"first_name": "Ryo",
		"last_name": "Ajeee",
		"email": "ryoaji27@gmail.com",
		"hire_date": "2023-06-27"
	}`
	result := model.CreateEmployeeResult{
		ID: 1,
	}
	testCases := []struct {
		Name                 string
		InitHandler          func(ctx echo.Context, s *mocks.EmployeeService) *Handler
		Json                 string
		ExpectedHttpCode     int
		ExpectedResponseBody model.ResponseBody
	}{
		{
			Name: "InvalidParams",
			InitHandler: func(ctx echo.Context, s *mocks.EmployeeService) *Handler {
				return NewHandler(s)
			},
			ExpectedHttpCode:     http.StatusBadRequest,
			ExpectedResponseBody: responseutil.CreateErrorResponse(pkgerror.ErrInvalidParams),
		},
		{
			Name: "ServiceError",
			InitHandler: func(ctx echo.Context, s *mocks.EmployeeService) *Handler {
				s.On("CreateEmployee", mock.Anything, mock.Anything).Return(nil, pkgerror.ErrSystemError)
				return NewHandler(s)
			},
			Json:                 validJson,
			ExpectedHttpCode:     http.StatusInternalServerError,
			ExpectedResponseBody: responseutil.CreateErrorResponse(pkgerror.ErrSystemError),
		},
		{
			Name: "Success",
			InitHandler: func(ctx echo.Context, s *mocks.EmployeeService) *Handler {
				s.On("CreateEmployee", mock.Anything, mock.Anything).Return(&result, pkgerror.NoError)
				return NewHandler(s)
			},
			Json:                 validJson,
			ExpectedHttpCode:     http.StatusOK,
			ExpectedResponseBody: responseutil.CreateSuccessResponse(&result, nil),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			v := validator.New()
			v.RegisterCustomTypeFunc(pkgvalidator.DecimalValidator, decimal.Decimal{})
			v.RegisterValidation("notblank", validators.NotBlank)
			e := echo.New()
			e.Validator = pkgvalidator.New(v)
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.Json))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			res := httptest.NewRecorder()
			c := e.NewContext(req, res)
			c.SetPath("/v1/employees")
			s := new(mocks.EmployeeService)
			h := tc.InitHandler(c, s)
			if assert.NoError(t, h.AddEmployee(c)) {
				assert.Equal(t, tc.ExpectedHttpCode, res.Code)
				expected := tc.ExpectedResponseBody
				jsonpath, err := jsonutil.NewJsonPath(res.Body.String())
				assert.Nil(t, err)
				assert.Equal(t, expected.Status, jsonpath.GetString("status"))
				assert.Equal(t, expected.Code, jsonpath.GetString("code"))
				assert.Equal(t, expected.ErrorMessage, jsonpath.GetStringPtr("error_message"))
				if expected.Data != nil {
					data := expected.Data.(*model.CreateEmployeeResult)
					assert.Equal(t, data.ID, jsonpath.GetInt("data.id"))
				}
			}
			s.AssertExpectations(t)
		})
	}
}
