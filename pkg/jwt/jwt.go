package jwt

import (
	"fmt"
	gojwt "github.com/dgrijalva/jwt-go"
)

func SignWithHS256(claims *gojwt.StandardClaims, secret []byte) (token string, err error) {
	return signWithHMAC(claims, gojwt.SigningMethodHS256, secret)
}

func VerifyWithHMAC(tokenString string, secret []byte) (*gojwt.StandardClaims, error) {
	claims := &gojwt.StandardClaims{}

	token, err := gojwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *gojwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*gojwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return secret, nil
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

func ExtractClaims(tokenString string) (*gojwt.StandardClaims, error) {
	claims := &gojwt.StandardClaims{}
	_, _, err := new(gojwt.Parser).ParseUnverified(tokenString, claims)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func signWithHMAC(claims *gojwt.StandardClaims, method *gojwt.SigningMethodHMAC, secret []byte) (token string, err error) {
	tokenClaims := gojwt.NewWithClaims(method, claims)
	token, err = tokenClaims.SignedString(secret)

	return
}
