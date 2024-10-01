package users

import (
	. "untitled/internal/musicstorage/ctl"
	. "untitled/internal/musicstorage/repo"
	. "untitled/internal/musicstorage/svc"
	. "untitled/internal/musicstorage/urls"

	"github.com/gin-gonic/gin"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"gorm.io/gorm"
)

type MusicStorageApp struct {
	ctl *MusicCtl
}

func NewMusicStorageApp(db *gorm.DB, tusdHandler *tusd.UnroutedHandler) *MusicStorageApp {
	musicRepo := NewTrackRepo(db)
	musicSvc := NewMusicSvc(musicRepo)

	return &MusicStorageApp{
		ctl: NewMusicCtl(musicSvc, tusdHandler),
	}
}

func (app *MusicStorageApp) Init(r *gin.Engine) {
	SetupMusicUrls(r, app.ctl)
}
