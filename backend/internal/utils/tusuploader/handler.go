package tusuploader

import (
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	tusd "github.com/tus/tusd/v2/pkg/handler"
)

const (
	maxFileSize = 50 * 1024 * 1024
	rGroup      = "/upload"
)

func InitTusUploader(r *gin.Engine, uploadDir string) {
	tusHandler, err := NewTusHandler(uploadDir)
	if err != nil {
		log.Fatalf("Failed to create TUS tusd: %v", err)
	}

	setupTusRoutes(r, tusHandler)
}

func NewTusHandler(uploadDir string) (*tusd.UnroutedHandler, error) {
	store := NewLocalStore(uploadDir)
	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	config := tusd.Config{
		BasePath:              rGroup,
		StoreComposer:         composer,
		MaxSize:               maxFileSize,
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

