package constant

import "errors"

type EmployeeColumn string

const (
	EmployeeColumnFirstName EmployeeColumn = "first_name"
	EmployeeColumnLastName  EmployeeColumn = "last_name"
	EmployeeColumnEmail     EmployeeColumn = "email"
	EmployeeColumnHireDate  EmployeeColumn = "hire_date"
)

var EmployeeColumns = []EmployeeColumn{
	EmployeeColumnFirstName,
	EmployeeColumnLastName,
	EmployeeColumnEmail,
	EmployeeColumnHireDate,
}

func ParseEmployeeColumnName(str string) (EmployeeColumn, error) {
	for _, t := range EmployeeColumns {
		if str == string(t) {
			return t, nil
		}
	}
	return "", errors.New(str)
}
