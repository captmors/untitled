package ctl

import (
	"net/http"
	"strings"
	"sync"
	. "untitled/internal/musicstorage/mdl"
	. "untitled/internal/musicstorage/svc"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/tus/tusd/v2/pkg/handler"
)

type MusicCtl struct {
	Svc        *MusicSvc
	tusHandler *handler.UnroutedHandler
}

func NewMusicCtl(musicSvc *MusicSvc, tusHandler *handler.UnroutedHandler) *MusicCtl {
	return &MusicCtl{Svc: musicSvc, tusHandler: tusHandler}
}

func (ctl *MusicCtl) UploadTrack(c *gin.Context) {
	var req UploadTrackRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uploadLength := c.Query("upload_length")
	if uploadLength == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Upload-Length is required"})
		return
	}

	c.Request.Header.Set("Tus-Resumable", "1.0.0")
	c.Request.Header.Set("Content-Type", "application/offset+octet-stream")
	c.Request.Header.Set("Upload-Length", uploadLength)
	c.Request.Header.Set("Upload-Offset", "0")

	authHeader := c.Request.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
		return
	}

	var wg sync.WaitGroup
	var uploadStatus int
	var uuid string
	var newTrackID int64

	// upload file resumably 
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("Starting file upload")
		ctl.tusHandler.PostFile(c.Writer, c.Request)
		uploadStatus = c.Writer.Status()
		uuid = c.Writer.Header().Get("Location")
		log.Info("File upload completed")
	}()

	// preadd it's metadata to the db
	wg.Add(1)
	go func() {
		defer wg.Done()
		userID, err := GetCurrentUserID(c)
		if err != nil {
			log.Error(err)
			return
		}

		newTrack := Track{
			UserID: int64(userID),
			Title:  c.Query("title"),
			Artist: c.Query("artist"),
			Format: c.Query("format"),
		}

		if err := ctl.Svc.Repo.Create(&newTrack); err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed to save track")
			return
		}

		newTrackID = newTrack.ID
		log.WithFields(log.Fields{"track_id": newTrackID}).Info("Track saved successfully")
	}()

	wg.Wait()

	// update physical ptr of file defined by uuid
	if uuid != "" {
		if err := ctl.Svc.Repo.UpdateTrackPtr(newTrackID, uuid); err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed to update track UUID ptr")
		} else {
			log.WithFields(log.Fields{"uuid": uuid}).Info("Track UUID ptr updated successfully")
		}
	}

	if uploadStatus != http.StatusCreated {
		if err := ctl.Svc.Repo.DeleteByID(newTrackID); err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed to delete track after upload failure")
		} else {
			log.WithFields(log.Fields{"track_id": newTrackID}).Info("Track deleted successfully after upload failure")
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed"})
		log.Error("Upload failed with status:", uploadStatus)
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Upload started"})
	log.Info("Upload process initiated")
}
