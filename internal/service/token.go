package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	key string
	exp time.Duration
}

func NewToken(key string, exp time.Duration) *Token {
	return &Token{
		key: key,
		exp: exp,
	}
}

func (t *Token) Generate(subject string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   subject,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.exp)),
	}
	
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(t.key))
}
