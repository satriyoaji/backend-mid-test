package service

import (
	"backend_test/model"
	"backend_test/pkg/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMain(m *testing.M) {
	err := config.LoadWithPath("./../configs/config-test.yml")
	if err != nil {
		log.Fatal("Load config error: ", err)
	}
	m.Run()
}

func createEchoContext(superadmin bool) echo.Context {
	claims := model.JwtClaims{}
	claims.Aud = []string{"backend-test"}
	claims.User.Superadmin = superadmin
	claims.User.ID = 1
	claims.User.Email = "user@gmail.com"
	e := echo.New()
	ctx := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
	ctx.Set("jwt_claims", &claims)
	return ctx
}
