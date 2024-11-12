package ctl

import (
	"net/http"
	"untitled/internal/users/mdl"
	"untitled/internal/users/svc"

	"github.com/gin-gonic/gin"
)

type AuthCtl struct {
	authSvc *svc.AuthSvc
}

func NewAuthCtl(authSvc *svc.AuthSvc) *AuthCtl {
	return &AuthCtl{
		authSvc: authSvc,
	}
}

func (ctl *AuthCtl) Login(c *gin.Context) {
    var req mdl.UserLoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := ctl.authSvc.AuthenticateByName(req.Name, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token, err := ctl.authSvc.GenerateToken(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, mdl.UserLoginResponse{Token: token})
}


func (ctl *AuthCtl) Register(c *gin.Context) {
    var req mdl.UserRegistrationRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := ctl.authSvc.Register(req.Name, req.Email, req.Password)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, mdl.UserRegistrationResponse{User: user})
}