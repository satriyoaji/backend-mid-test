package model

import "time"

type GetEmployeesFilter struct {
	FirstName   string `query:"first_name"`
	LastName    string `query:"last_name"`
	ID          *int   `query:"id"`
	PageRequest PageRequest
}

type GetEmployeesResult struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	HireDate  time.Time `json:"hire_date"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
}

type CreateEmployeeRequest struct {
	FirstName string `json:"first_name" validate:"required,notblank,min=3,max=60"`
	LastName  string `json:"last_name" validate:"required,notblank,min=3,max=60"`
	Email     string `json:"email" validate:"required,notblank,email,min=3,max=60"`
	HireDate  string `json:"hire_date" validate:"required,notblank"`
}

type CreateEmployeeResult struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	HireDate  time.Time `json:"hire_date"`
}

type GetEmployeeByIDRequest struct {
	EmployeeID int `param:"id" validate:"required"`
}

type DeleteEmployeeByIDRequest struct {
	EmployeeID int `param:"id" validate:"required"`
}

type GetEmployeeByIDResult struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	HireDate  time.Time `json:"hire_date"`
}

type EditEmployeeRequest struct {
	EmployeeID int `param:"id" validate:"required"` // Path variable

	FirstName string `json:"first_name" validate:"required,notblank,min=3,max=60"`
	LastName  string `json:"last_name" validate:"required,notblank,min=3,max=60"`
	Email     string `json:"email" validate:"required,notblank,email,min=3,max=60"`
	HireDate  string `json:"hire_date" validate:"required,notblank"`
}

type EditEmployeeResult struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	HireDate  time.Time `json:"hire_date"`
}
