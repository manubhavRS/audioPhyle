package models

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
	ID             string    `db:"user_id" json:"userID"`
	Exp            time.Time `db:"exp" json:"exp"`
	StandardClaims jwt.StandardClaims
}

func (c Claims) Valid() error {
	return c.StandardClaims.Valid()
}
