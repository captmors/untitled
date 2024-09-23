package urls

import (
	"untitled/internal/users/ctl"

	"github.com/gin-gonic/gin"
)

func SetupUserUrls(r *gin.Engine, userCtl *ctl.UserCtl) {
	userRoutes := r.Group("/users")
	{
		userRoutes.GET("/:id", userCtl.GetUserByID)
		userRoutes.POST("", userCtl.CreateUser)
	}
}
