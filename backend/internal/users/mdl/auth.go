package mdl

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type UserLoginRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type UserRegistrationRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRegistrationResponse struct {
	User *User `json:"user"`
}
