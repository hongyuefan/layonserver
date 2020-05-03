package token

import (
	"errors"
	"fmt"

	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	Expire_Time = 120
	Secret_Key  = "this is Layon !"
)

func TokenGenerate(userId int64) (token string, err error) {

	ojwt := jwt.New(jwt.SigningMethodHS256)

	ojwt.Claims = &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * Expire_Time).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    fmt.Sprintf("%v", userId),
	}

	return ojwt.SignedString([]byte(Secret_Key))
}

func TokenValidate(token string) (userId string, err error) {

	ojwt, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret_Key), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := ojwt.Claims.(*jwt.StandardClaims); ok && ojwt.Valid {
		return claims.Issuer, nil
	}

	return "", errors.New("token valide failse,please login again")

}
