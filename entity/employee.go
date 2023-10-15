package entity

import (
	"time"
)

type Employee struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	FirstName string
	LastName  string
	Email     string
	HireDate  time.Time
}

func (Employee) TableName() string {
	return "employees"
}
