package tusuploader

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	tusd "github.com/tus/tusd/v2/pkg/handler"
)

const (
	rGroup = "/upload"
)

type TusHandlerCfg struct {
	MaxFileSize int64
	UploadDir   string
}

func InitTusUploader(r *gin.Engine, cfg TusHandlerCfg) {
	tusHandler, err := NewTusHandler(cfg)
	if err != nil {
		log.Fatalf("Failed to create TUS tusd: %v", err)
	}

	setupTusRoutes(r, tusHandler)
}

func NewTusHandler(cfg TusHandlerCfg) (*tusd.UnroutedHandler, error) {
	store := NewLocalStore(cfg.UploadDir)
	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	config := tusd.Config{
		BasePath:              rGroup,
		StoreComposer:         composer,
		MaxSize:               cfg.MaxFileSize,
		NotifyCompleteUploads: true,
	}
	tusHandler, err := tusd.NewUnroutedHandler(config)

	if err != nil {
		return nil, err
	}

	return tusHandler, nil
}

func setupTusRoutes(r *gin.Engine, tusHandler *tusd.UnroutedHandler) {
	tusRoutes := r.Group(rGroup)
	{
		tusRoutes.POST("", gin.WrapF(tusHandler.PostFile))
		tusRoutes.HEAD("/:id", gin.WrapF(tusHandler.HeadFile))
		tusRoutes.PATCH("/:id", gin.WrapF(tusHandler.PatchFile))
	}
}
