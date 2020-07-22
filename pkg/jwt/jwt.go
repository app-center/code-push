package jwt

import (
	"fmt"
	gojwt "github.com/dgrijalva/jwt-go"
	"time"
)

func NewJwt(fns ...func(*options)) *jwt {
	jwtOptions := &options{
		Secret:   "",
		Issuer:   "",
		Lifetime: 0,
	}

	for _, fn := range fns {
		fn(jwtOptions)
	}

	return &jwt{options: jwtOptions}
}

type options struct {
	Secret   string
	Issuer   string
	Lifetime time.Duration
}

type jwt struct {
	options *options
}

func (j *jwt) SignToken(subject string) (token string, err error) {
	expireTime := time.Now().Add(j.options.Lifetime)

	claims := gojwt.StandardClaims{
		ExpiresAt: expireTime.Unix(),
		Issuer:    j.options.Issuer,
		Subject:   subject,
	}

	tokenClaims := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	token, err = tokenClaims.SignedString(j.options.Secret)

	return
}

func (j *jwt) VerifyToken(tokenString string) (*gojwt.StandardClaims, error) {
	claims := &gojwt.StandardClaims{}

	token, err := gojwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *gojwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*gojwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return j.options.Secret, nil
		})

	if token != nil {
		if token.Valid {
			return claims, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
