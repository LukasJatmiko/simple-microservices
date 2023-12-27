package middlewares

import (
	"crypto/rsa"
	"fmt"
	"net/http"

	"github.com/LukasJatmiko/simple-microservices/data-processor/packages/api"
	"github.com/LukasJatmiko/simple-microservices/data-processor/packages/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/ssh"
)

type AuthJWT struct {
	Options *auth.Options
}

func NewAuthJWT(opts *auth.Options) *AuthJWT {
	return &AuthJWT{Options: opts}
}

func (a *AuthJWT) GetOptions() *auth.Options {
	return a.Options
}

func (a *AuthJWT) ValidateAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if tokenString, e := GetAuthToken(c); e != nil {
			return c.JSON(http.StatusOK, api.APIResponsePayload{Status: http.StatusUnauthorized, Message: "unauthorized", Data: echo.Map{}})
		} else {
			if token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				} else {
					if pubk, _, _, _, e := ssh.ParseAuthorizedKey(a.GetOptions().AuthJWTPublicKey); e != nil {
						return nil, fmt.Errorf("error parsing public key: %v", e)
					} else {
						if sshpubkey, e := ssh.ParsePublicKey(pubk.Marshal()); e != nil {
							return nil, fmt.Errorf("error parsing public key: %v", e)
						} else {
							pubCrypto := sshpubkey.(ssh.CryptoPublicKey).CryptoPublicKey()
							publicKey := pubCrypto.(*rsa.PublicKey)
							return publicKey, nil
						}
					}
				}
			}); err != nil {
				return c.JSON(http.StatusOK, api.APIResponsePayload{Status: http.StatusUnauthorized, Message: "unauthorized", Data: echo.Map{}})
			} else {
				if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					return next(c)
				} else {
					return c.JSON(http.StatusOK, api.APIResponsePayload{Status: http.StatusUnauthorized, Message: "unauthorized", Data: echo.Map{}})
				}
			}
		}
	}
}
