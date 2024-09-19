package mdl

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

type UserLoginRequest struct {
    Email    string `json:"email" binding:"required"`
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
