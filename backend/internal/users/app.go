package users

import (
	"untitled/internal/users/ctl"
	"untitled/internal/users/mw"
	"untitled/internal/users/repo"
	"untitled/internal/users/svc"
	"untitled/internal/users/urls"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type UserApp struct {
	ctl          *ctl.UserCtl
	authCtl      *ctl.AuthCtl
	authMWConfig mw.AuthMWConfig
}

func NewUserApp(db *gorm.DB, jwtKey []byte) *UserApp {
	userRepo := repo.NewUserRepo(db)
	userSvc := svc.NewUserSvc(userRepo)
	authSvc := svc.NewAuthSvc(userRepo, jwtKey)

	authMWConfig := mw.AuthMWConfig{
		JwtKey: jwtKey,
		Claims: func() jwt.Claims {
			return &jwt.RegisteredClaims{}
		},
	}

	return &UserApp{
		ctl:          ctl.NewUserCtl(userSvc),
		authCtl:      ctl.NewAuthCtl(authSvc),
		authMWConfig: authMWConfig,
	}
}

func (app *UserApp) Init(r *gin.Engine) {
	r.Use(mw.AuthMW(app.authMWConfig))

	urls.SetupUserUrls(r, app.ctl)
	urls.SetupAuthUrls(r, app.authCtl)
}
