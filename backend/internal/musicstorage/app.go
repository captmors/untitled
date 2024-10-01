package users

import (
	. "untitled/internal/musicstorage/ctl"
	. "untitled/internal/musicstorage/repo"
	. "untitled/internal/musicstorage/svc"
	. "untitled/internal/musicstorage/urls"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	tusd "github.com/tus/tusd/v2/pkg/handler"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type MusicStorageApp struct {
	ctl *MusicCtl
}

func NewMusicStorageApp(mongoDB *mongo.Client, pgDB *gorm.DB, esClient *es.Client, tusdHandler *tusd.UnroutedHandler) *MusicStorageApp {
    mongoRepo := NewMongoRepo(mongoDB, "musicstorage", "tracks") 
    pgRepo := NewPgRepo(pgDB) 

    musicSvc := NewMusicSvc(pgRepo, mongoRepo)

    return &MusicStorageApp{
        ctl: NewMusicCtl(musicSvc, tusdHandler),
    }
}


func (app *MusicStorageApp) Init(r *gin.Engine) {
	SetupMusicUrls(r, app.ctl)
}
