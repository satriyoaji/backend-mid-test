package handler

import (
	mocks "backend_test/mocks/service"
	"backend_test/pkg/config"
	"github.com/labstack/echo/v4"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := config.LoadWithPath("./../../../configs/config-test.yml")
	if err != nil {
		log.Fatal("Load config error: ", err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestRegisterHandlers(t *testing.T) {
	h := NewHandler(
		&mocks.EmployeeService{},
	)
	RegisterHandlers(echo.New(), h)
}
