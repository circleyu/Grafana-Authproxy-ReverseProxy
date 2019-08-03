package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var hmacSecret []byte

type myCustomClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

// CreateToken 產生token
func CreateToken(userName string) (string, error) {

	claims := myCustomClaims{
		userName,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(15)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(hmacSecret)
}
