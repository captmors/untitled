package users

import (
	"untitled/internal/users/ctl"
	"untitled/internal/users/repo"
	"untitled/internal/users/svc"
	"untitled/internal/users/urls"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserApp struct {
	ctl *ctl.UserCtl
}

func NewUserApp(db *gorm.DB) *UserApp {
	return &UserApp{
		ctl: ctl.NewUserCtl(svc.NewUserSvc(repo.NewUserRepo(db))),
	}
}

func (app *UserApp) Init(r *gin.Engine) {
	urls.SetupUrls(r, app.ctl)
}
