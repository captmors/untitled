package urls

import (
	"untitled/internal/musicstorage/ctl"

	"github.com/gin-gonic/gin"
)

func SetupMusicUrls(r *gin.Engine, musicCtl *ctl.MusicCtl) {
	trackRoutes := r.Group("/tracks")
	{
		trackRoutes.POST("/upload", musicCtl.UploadTrack)
		trackRoutes.GET("/", musicCtl.ListTracks)
		trackRoutes.GET("/:id", musicCtl.GetTrackByID)
		trackRoutes.PUT("/:id", musicCtl.UpdateTrackMetadata)
		trackRoutes.DELETE("/:id", musicCtl.RemoveTrack)
	}
}
