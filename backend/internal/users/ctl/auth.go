package ctl

import (
    "net/http"
    "untitled/internal/users/mdl"
    "untitled/internal/users/svc"

    "github.com/gin-gonic/gin"
)

type AuthCtl struct {
    authService *svc.AuthService
}

func NewAuthCtl(authService *svc.AuthService) *AuthCtl {
    return &AuthCtl{
        authService: authService,
    }
}

func (ctl *AuthCtl) Login(c *gin.Context) {
    var req mdl.UserLoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := ctl.authService.Authenticate(req.Email, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token, err := ctl.authService.GenerateToken(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, mdl.UserLoginResponse{Token: token})
}
