package entity

import (
	 "github.com/golang-jwt/jwt/v5"
	 "os"

)
	

var JwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type UserJwt struct {
	UserId int64 `json:"userId"`
	Username string `json:"username"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claimes struct {
	UserId int64 `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
