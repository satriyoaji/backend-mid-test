package main

import (
	"backend_test/cmd/app/handler"
	"backend_test/pkg/config"
	"backend_test/pkg/db"
	pkgvalidator "backend_test/pkg/validator"
	"backend_test/repository"
	"backend_test/service"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/shopspring/decimal"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}
	if !config.Data.IsEnvProduction() {
		log.SetLevel(log.DEBUG)
	}

	dbh := db.Init()
	db.Migrate(dbh)

	repo := repository.Default(dbh)

	employeeService := service.NewEmployeeService(repo)

	h := handler.NewHandler(employeeService)

	v := validator.New()
	v.RegisterCustomTypeFunc(pkgvalidator.DecimalValidator, decimal.Decimal{})
	v.RegisterValidation("notblank", validators.NotBlank)

	e := echo.New()
	e.Validator = pkgvalidator.New(v)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	handler.RegisterHandlers(e, h)
	err = e.Start(config.Data.Port)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
