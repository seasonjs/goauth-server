package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
)

// NewJwtManager JWT token构造函数
func NewJwtManager(signedKeyID string, signedKey string, method jwt.SigningMethod) *manage.Manager {
	m := manage.NewManager()
	m.MapAuthorizeGenerate(generates.NewAuthorizeGenerate())
	//m.MapAccessGenerate(generatesNewJWTAccessGenerate(signedKeyID, []byte(signedKey), method))
	return m
}

// ParseJwtToken 解析与验证JwtToken
func ParseJwtToken(access string, signedKey string) *generates.JWTAccessClaims {
	//解析验证Token
	token, err := jwt.ParseWithClaims(access, &generates.JWTAccessClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("parse error")
		}
		return []byte(signedKey), nil
	})
	if err != nil {
		panic(err)
	}

	claims, ok := token.Claims.(*generates.JWTAccessClaims)
	if !ok || !token.Valid {
		panic("invalid token")
	}
	return claims
}
