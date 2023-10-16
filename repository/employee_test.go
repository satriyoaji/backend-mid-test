package repository

import (
	"backend_test/entity"
	"backend_test/model"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindEmployees(t *testing.T) {
	shops, err := repo.FindEmployees(context.Background(), model.GetEmployeesFilter{})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(shops))
}

func TestFindEmployeeByID(t *testing.T) {
	person, err := repo.FindEmployeeByID(context.Background(), uint(1))
	assert.Nil(t, err)
	assert.Equal(t, uint(1), person.ID)
	assert.Equal(t, "First Employee 1", person.FirstName)
	assert.Equal(t, "Last 1", person.LastName)
}

func TestUpdateEmployee(t *testing.T) {
	person := entity.Employee{
		ID:        1,
		FirstName: "First name",
		LastName:  "Last name",
	}
	err := repo.UpdateEmployee(context.Background(), &person)
	assert.Nil(t, err)
	result, err := repo.FindEmployeeByID(context.Background(), 1)
	assert.Nil(t, err)
	assert.Equal(t, person.ID, result.ID)
	assert.Equal(t, person.FirstName, result.FirstName)
	assert.Equal(t, person.LastName, result.LastName)
	resetData()
}
