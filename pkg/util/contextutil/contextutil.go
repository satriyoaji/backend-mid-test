package contextutil

import (
	"backend_test/model"

	"github.com/labstack/echo/v4"
)

func GetJwtClaims(ctx echo.Context) (claims *model.JwtClaims) {
	c := ctx.Get("jwt_claims")
	if c == nil {
		return nil
	}
	return c.(*model.JwtClaims)
}

func GetUserEmail(ctx echo.Context) *string {
	claims := GetJwtClaims(ctx)
	if claims == nil {
		return nil
	}
	return &claims.User.Email
}

func GetAppIDsFromJwt(ctx echo.Context) *[]int {
	appIDs := []int{}
	claims := GetJwtClaims(ctx)
	if claims == nil {
		return &appIDs
	}
	for _, role := range claims.User.Roles {
		appIDs = append(appIDs, role.App.ID)
	}
	return &appIDs
}

func IsSuperadmin(ctx echo.Context) bool {
	claims := GetJwtClaims(ctx)
	if claims == nil {
		return false
	}
	return claims.User.Superadmin
}

func AppIDsContain(ctx echo.Context, appID int) bool {
	claims := GetJwtClaims(ctx)
	if claims == nil {
		return false
	}
	for _, role := range claims.User.Roles {
		if role.App.ID == appID {
			return true
		}
	}
	return false
}
