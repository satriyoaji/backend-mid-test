package repository

import (
	"backend_test/constant"
	"backend_test/entity"
	"backend_test/model"
	"context"

	"gorm.io/gorm"
)

func whereEmployeeFirstNameContains(name string, alias string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if name == "" {
			return db
		}
		sql := withAlias(string(constant.EmployeeColumnFirstName), alias) + " ilike ?"
		return db.Where(sql, withPercentAround(name))
	}
}
func whereEmployeeLastNameContains(name string, alias string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if name == "" {
			return db
		}
		sql := withAlias(string(constant.EmployeeColumnLastName), alias) + " ilike ?"
		return db.Where(sql, withPercentAround(name))
	}
}

func whereEmployeeIDIn(ids []int, alias string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(ids) == 0 {
			return db
		}
		sql := withAlias("id", alias) + " in ?"
		return db.Where(sql, ids)
	}
}

func (d DefaultRepository) FindEmployees(ctx context.Context, filter model.GetEmployeesFilter) ([]entity.Employee, error) {
	employeeIDs := []int{}
	if filter.ID != nil {
		employeeIDs = append(employeeIDs, *filter.ID)
	}
	shops := []entity.Employee{}
	err := d.handler.Tx.WithContext(ctx).
		Scopes(
			whereEmployeeFirstNameContains(filter.FirstName, ""),
			whereEmployeeLastNameContains(filter.LastName, ""),
			whereEmployeeIDIn(employeeIDs, ""),
			paginate(filter.PageRequest.PageNum, filter.PageRequest.PageSize)).
		Order("created_at desc").Find(&shops).Error
	return shops, err
}

func (d DefaultRepository) FindAllEmployees(ctx context.Context, filter model.GetEmployeesFilter) ([]entity.Employee, error) {
	employeeIDs := []int{}
	if filter.ID != nil {
		employeeIDs = append(employeeIDs, *filter.ID)
	}
	shops := []entity.Employee{}
	err := d.handler.Tx.WithContext(ctx).
		Scopes(
			whereEmployeeFirstNameContains(filter.FirstName, ""),
			whereEmployeeLastNameContains(filter.LastName, ""),
			whereEmployeeIDIn(employeeIDs, "")).
		Order("created_at desc").Find(&shops).Error
	return shops, err
}

func (d DefaultRepository) CountEmployees(ctx context.Context, filter model.GetEmployeesFilter) (int, error) {
	employeeIDs := []int{}
	if filter.ID != nil {
		employeeIDs = append(employeeIDs, *filter.ID)
	}
	var count int64
	err := d.handler.Tx.WithContext(ctx).Model(&entity.Employee{}).
		Scopes(
			whereEmployeeFirstNameContains(filter.FirstName, ""),
			whereEmployeeLastNameContains(filter.LastName, ""),
			whereEmployeeIDIn(employeeIDs, "")).
		Count(&count).Error
	return int(count), err
}

func (d DefaultRepository) CreateEmployee(ctx context.Context, employee *entity.Employee) error {
	return d.handler.Tx.Create(employee).Error
}

func (d DefaultRepository) FindEmployeeByID(ctx context.Context, id uint) (entity.Employee, error) {
	employee := entity.Employee{}
	err := d.handler.Tx.WithContext(ctx).Where("id=?", id).First(&employee).Error
	return employee, err
}

func (d DefaultRepository) FindEmployeeByEmail(ctx context.Context, email string) (entity.Employee, error) {
	employee := entity.Employee{}
	err := d.handler.Tx.WithContext(ctx).Where("email=?", email).First(&employee).Error
	return employee, err
}

func (d DefaultRepository) UpdateEmployee(ctx context.Context, employee *entity.Employee) error {
	return d.handler.Tx.WithContext(ctx).Save(employee).Error
}

func (d DefaultRepository) DeleteEmployee(ctx context.Context, id uint) error {
	return d.handler.Tx.WithContext(ctx).Delete(&entity.Employee{}, id).Error
}
