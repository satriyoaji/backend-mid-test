package repository

import (
	"backend_test/entity"
	"backend_test/model"
	"backend_test/pkg/db"
	"context"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"strings"
)

type Repository interface {
	TxBegin() error
	TxCommit() error
	TxRollback() error

	// Employee
	FindEmployees(ctx context.Context, filter model.GetEmployeesFilter) ([]entity.Employee, error)
	FindAllEmployees(ctx context.Context, filter model.GetEmployeesFilter) ([]entity.Employee, error)
	CountEmployees(ctx context.Context, filter model.GetEmployeesFilter) (int, error)
	CreateEmployee(ctx context.Context, merchant *entity.Employee) error
	FindEmployeeByID(ctx context.Context, id uint) (entity.Employee, error)
	UpdateEmployee(ctx context.Context, merchant *entity.Employee) error
	DeleteEmployee(ctx context.Context, id uint) error
}

type DefaultRepository struct {
	handler *db.Handler
}

func Default(handler *db.Handler) *DefaultRepository {
	return &DefaultRepository{
		handler: handler,
	}
}

func (d DefaultRepository) TxBegin() error {
	log.Debug("Start db transaction")
	d.handler.Tx = d.handler.DB.Begin()
	return d.handler.Tx.Error
}

func (d DefaultRepository) TxCommit() error {
	log.Debug("Commit db transaction")
	err := d.handler.Tx.Commit().Error
	d.handler.Tx = d.handler.DB
	return err
}

func (d DefaultRepository) TxRollback() error {
	log.Debug("Rollback db transaction")
	err := d.handler.Tx.Rollback().Error
	d.handler.Tx = d.handler.DB
	return err
}

func paginate(pageNum, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNum <= 0 {
			pageNum = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}
		offset := (pageNum - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func withAlias(column, alias string) string {
	if alias == "" {
		return column
	}
	return alias + "." + column
}

func withPercentAround(val string) string {
	return "%" + val + "%"
}

func withPercentAfter(val string) string {
	return val + "%"
}

func withPercentBefore(val string) string {
	return "%" + val
}

func getSortDir(sortDir string) string {
	if strings.ToLower(sortDir) == "asc" {
		return strings.ToLower(sortDir)
	}
	return "desc"
}
