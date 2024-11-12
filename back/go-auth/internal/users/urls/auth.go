package urls

import (
	"untitled/internal/users/ctl"

	"github.com/gin-gonic/gin"
)

func SetupAuthUrls(r *gin.Engine, authCtl *ctl.AuthCtl) {
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/login", authCtl.Login)
		authRoutes.POST("/register", authCtl.Register)
	}
}
