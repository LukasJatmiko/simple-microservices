package auth

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/LukasJatmiko/simple-microservices/data-processor/driver"
	"github.com/LukasJatmiko/simple-microservices/data-processor/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type UserClaims struct {
	UserID    int32  `json:"user_id" gorm:"column:user_id"`
	UserEmail string `json:"user_email" gorm:"column:user_email"`
	UserName  string `json:"user_name" gorm:"column:user_name"`
	jwt.RegisteredClaims
}

type AuthJWT struct {
	JWTSigningKey []byte
	JWTExpiration time.Duration
	JWTIssuer     string
	Driver        driver.Driver
}

func (auth *AuthJWT) Login(cred Credential) (*Data, error) {
	claims := new(UserClaims)
	userdata := make(map[string]interface{})
	switch auth.Driver.GetWrapperInstance().(type) {
	case *gorm.DB:
		{
			sql := auth.Driver.GetWrapperInstance().(*gorm.DB)
			query := `
			select
				u.id as user_id,
				u.name as user_name,
				u.email as user_email,
				u.password as user_password
			from users u
			where email = ?;
			`
			if result := sql.Raw(query, cred.Email); result.Error != nil {
				return nil, fmt.Errorf("unknown user")
			} else {
				if e := result.Scan(userdata).Error; e != nil || userdata["user_id"] == nil {
					return nil, fmt.Errorf("unknown user")
				} else {
					if e := utils.Map(userdata).Decode(claims); e != nil {
						return nil, fmt.Errorf("unknown user")
					} else {
						if e := utils.BcryptCheckPassword(cred.Password, userdata["user_password"].(string)); e != nil {
							return nil, fmt.Errorf("incorrect password")
						}
					}
				}
			}
		}
	default:
		{
			//to do
		}
	}
	if id, e := uuid.NewRandom(); e != nil {
		return nil, fmt.Errorf("error while creating uuid (%v)", e)
	} else {
		claims.ID = id.String()
	}
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(auth.JWTExpiration))
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.NotBefore = jwt.NewNumericDate(time.Now())
	claims.Issuer = auth.JWTIssuer
	claims.Subject = claims.UserName
	claims.Audience = []string{auth.JWTIssuer}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	block, _ := pem.Decode(auth.JWTSigningKey)
	privateKey, e := x509.ParsePKCS1PrivateKey(block.Bytes)
	if e != nil {
		return nil, fmt.Errorf("error while parsing pivate key %v", e)
	}
	if tokenString, err := token.SignedString(privateKey); err != nil {
		return nil, fmt.Errorf("error while signing token %v", err)
	} else {
		return &Data{Token: tokenString}, nil
	}
}
