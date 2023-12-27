package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/LukasJatmiko/simple-microservices/data-processor/constants"
	"github.com/LukasJatmiko/simple-microservices/data-processor/driver"
	"github.com/LukasJatmiko/simple-microservices/data-processor/packages/api"
	"github.com/LukasJatmiko/simple-microservices/data-processor/types"
	"github.com/labstack/echo/v4"
)

type Options struct {
	AuthType          types.AuthType
	AuthTokens        []string
	AuthJWTPrivateKey []byte
	AuthJWTPublicKey  []byte
	AUthJWTExpiration time.Duration
	AuthJWTIssuer     string
	ActiveMiddlewares []string
}

type ContextKey string

type AuthRequestPayload struct {
	Credential Credential `json:"data"`
}

type Credential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Data struct {
	Token string `json:"token"`
}

type AuthHandler interface {
	Login(cred Credential) (*Data, error)
}

func NewAuthHandler(opts *Options, driver driver.Driver) AuthHandler {
	switch opts.AuthType {
	case constants.AuthTypeJWT:
		return &AuthJWT{
			JWTSigningKey: opts.AuthJWTPrivateKey,
			JWTExpiration: opts.AUthJWTExpiration,
			JWTIssuer:     opts.AuthJWTIssuer,
			Driver:        driver,
		}
	default:
		return &AuthStatic{AuthTokens: opts.AuthTokens}
	}
}

func Mount(opts *Options, handler AuthHandler, app *echo.Group) {
	app.POST("/login", Login(handler))
}

func Login(auth AuthHandler) func(c echo.Context) error {
	return func(c echo.Context) error {
		payload := new(AuthRequestPayload)
		if e := c.Bind(payload); e != nil {
			return c.JSON(http.StatusOK, api.APIResponsePayload{Status: http.StatusUnauthorized, Message: "unauthorized", Data: echo.Map{}})
		} else {
			if data, e := auth.Login(payload.Credential); e != nil {
				fmt.Println(e)
				return c.JSON(http.StatusOK, api.APIResponsePayload{Status: http.StatusUnauthorized, Message: "unauthorized", Data: echo.Map{}})
			} else {
				return c.JSON(http.StatusOK, api.APIResponsePayload{Status: http.StatusOK, Message: "ok", Data: data})
			}
		}
	}
}
