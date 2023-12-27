package auth

import "strings"

type AuthStatic struct {
	AuthTokens []string
}

func (a *AuthStatic) Login(cred Credential) (*Data, error) {
	return &Data{Token: strings.Join(a.AuthTokens, ";")}, nil
}
