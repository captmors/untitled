package ctl

import (
	"net/http"
	"strconv"
	"untitled/internal/users/mdl"
	"untitled/internal/users/svc"

	"github.com/gin-gonic/gin"
)

type UserCtl struct {
	userSvc *svc.UserSvc
}

func NewUserCtl(svc *svc.UserSvc) *UserCtl {
	return &UserCtl{
		userSvc: svc,
	}
}

func (ctl *UserCtl) CreateUser(c *gin.Context) {
	var user mdl.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdUser, err := ctl.userSvc.CreateUser(user.Name, user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, createdUser)
}

func (ctl *UserCtl) GetUserByID(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }
    user, err := ctl.userSvc.GetUserByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }
    c.JSON(http.StatusOK, user)
}