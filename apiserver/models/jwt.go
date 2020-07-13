package models

import (
	"apiserver/utils"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var secretKey = "$%^dfvthe2345gcvbjmERddOPWWKqZAP<q%^)+/f"

type UserAuthInfo struct {
	UserId     int              `json:"userId,string"`
	UserName   string           `json:"userName"`
	UserType   utils.UserType   `json:"userType,string"`
	UserStatus utils.UserStatus `json:"userStatus,string"`
	Gender     int              `json:"gender,string"`
	Birthday   *time.Time       `json:"birthDay"`
	Tel        string           `json:"tel"`
	Email      string           `json:"email"`
	Addr       string           `json:"addr"`
	Remark     string           `json:"remark"`
	Avatar     NullString       `json:"avatar"`
}

type AuthCustomClaims struct {
	UserAuthInfo
	jwt.StandardClaims
}

func ParseJwtAuthToken(tokenString string) *AuthCustomClaims {

	token, err := jwt.ParseWithClaims(tokenString, &AuthCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		utils.Logger.Error(err)
		return nil
	}

	if claims, ok := token.Claims.(*AuthCustomClaims); ok && token.Valid {
		return claims
	} else {
		return nil
	}
}

func GenJwtAuthToken(payload *AuthCustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenStr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		utils.Logger.Error(err)
		return "", err
	}
	return tokenStr, nil
}
