package middlewares

import (
	"fmt"
	"strings"

	"github.com/LukasJatmiko/simple-microservices/data-processor/constants"
	"github.com/LukasJatmiko/simple-microservices/data-processor/packages/auth"
	"github.com/labstack/echo/v4"
)

type Auth interface {
	GetOptions() *auth.Options
	ValidateAuth(next echo.HandlerFunc) echo.HandlerFunc
}

func NewAuthMiddleware(opts *auth.Options) (Auth, error) {
	switch strings.ToUpper(string(opts.AuthType)) {
	case string(constants.AuthTypeStatic):
		{
			return NewAuthStatic(opts), nil
		}
	case string(constants.AuthTypeJWT):
		{
			return NewAuthJWT(opts), nil
		}
	default:
		return nil, fmt.Errorf("unsupported auth type")
	}
}

func Mount(app *echo.Group, amw Auth) {
	app.Use(amw.ValidateAuth)
}

func GetAuthToken(c echo.Context) (string, error) {
	header := c.Request().Header
	authHeader := header.Get("Authorization")
	if len(authHeader) <= 8 {
		return "", fmt.Errorf("unauthorized")
	} else {
		return authHeader[7:], nil
	}
}
