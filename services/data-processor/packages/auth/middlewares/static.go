package middlewares

import (
	"net/http"

	"github.com/LukasJatmiko/simple-microservices/data-processor/packages/api"
	"github.com/LukasJatmiko/simple-microservices/data-processor/packages/auth"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slices"
)

type AuthStatic struct {
	Options *auth.Options
}

func NewAuthStatic(opts *auth.Options) *AuthStatic {
	return &AuthStatic{Options: opts}
}

func (a *AuthStatic) GetOptions() *auth.Options {
	return a.Options
}

func (a *AuthStatic) ValidateAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if token, e := GetAuthToken(c); e != nil {
			return c.JSON(http.StatusOK, api.APIResponsePayload{Status: http.StatusUnauthorized, Message: "unauthorized", Data: echo.Map{}})
		} else {
			idx := slices.IndexFunc(a.Options.AuthTokens, func(t string) bool { return t == token })
			if idx < 0 {
				return c.JSON(http.StatusOK, api.APIResponsePayload{Status: http.StatusUnauthorized, Message: "unauthorized", Data: echo.Map{}})
			} else {
				return next(c)
			}
		}
	}
}
