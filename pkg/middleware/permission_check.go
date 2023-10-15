package middleware

import (
	"backend_test/model"
	"backend_test/pkg/config"
	pkgerror "backend_test/pkg/error"
	"backend_test/pkg/util/responseutil"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/exp/slices"
)

// Path and permissions mapping
var mapping = map[string][]string{
	http.MethodGet + "/merchants":             {"read_merchants", "create_promotions"},
	http.MethodGet + "/merchants/:merchantId": {"read_merchants"},
	http.MethodPut + "/merchants/:merchantId": {"update_merchants"},
	http.MethodPost + "/merchants":            {"create_merchants"},
	http.MethodGet + "/products":              {"read_products", "create_promotions"},
	http.MethodGet + "/products/:productId":   {"read_products"},
	http.MethodPost + "/products":             {"create_products"},
	http.MethodPut + "/products/:productId":   {"update_products"},
	http.MethodPost + "/upload-image": {
		"create_merchants",
		"update_merchants",
		"create_products",
		"update_products",
	},
	http.MethodGet + "/campaigns":                {"read_campaigns"},
	http.MethodGet + "/campaigns/:campaignId":    {"read_campaigns"},
	http.MethodPost + "/campaigns":               {"create_campaigns"},
	http.MethodPut + "/campaigns/:campaignId":    {"update_campaigns"},
	http.MethodDelete + "/campaigns/:campaignId": {"delete_campaigns"},
	http.MethodGet + "/orders":                   {"read_orders"},
	http.MethodGet + "/orders/:orderId":          {"read_orders"},
	http.MethodGet + "/promotions":               {"read_promotions"},
	http.MethodGet + "/promotions/:promotionId":  {"read_promotions"},
	http.MethodPost + "/promotions":              {"create_promotions"},
	http.MethodPut + "/promotions/:promotionId":  {"update_promotions"},
}

func withAppName(names ...string) []string {
	permissions := []string{}
	for _, n := range names {
		permissions = append(permissions, config.Data.AppCode+":"+n)
	}
	return permissions
}

func validateUserPermission(method, path string, claims *model.JwtClaims) pkgerror.CustomError {
	key := method + path
	permissions, keyExists := mapping[key]
	if !keyExists || len(permissions) == 0 {
		// If the mapping does not contain path or the permissions is empty,
		// the request does not need any specific permission
		return pkgerror.NoError
	}
	prefixedPermissions := withAppName(permissions...)
	for _, role := range claims.User.Roles {
		for _, p := range prefixedPermissions {
			if slices.Contains(role.Permissions, p) {
				return pkgerror.NoError
			}
		}
	}
	log.Errorf("None of the required permissions %v found in the JWT", permissions)
	return pkgerror.ErrForbiddenRequest.WithError(fmt.Errorf("missing permissions: %v", permissions))
}

func PermissionCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		c := ctx.Get("jwt_claims")
		if c == nil {
			return next(ctx)
		}
		claims := c.(*model.JwtClaims)
		if claims.User.Superadmin {
			return next(ctx)
		}
		if err := validateUserPermission(ctx.Request().Method, ctx.Path(), claims); !err.IsNoError() {
			return responseutil.SendErrorResponse(ctx, err)
		}
		return next(ctx)
	}
}
