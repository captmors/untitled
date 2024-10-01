package ctl

import (
	"net/http"
	"strconv"
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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// prepare request to tusd server
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
		uuid = GetUUIDFromLocation(c.Writer.Header().Get("Location"))
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
			UserID:   int64(userID),
			Title:    req.Title,
			Artist:   req.Artist,
			Format:   req.Format,
			Bitrate:  req.Bitrate,
			Duration: req.Duration,
			Genre:    req.Genre,
		}

		if err := ctl.Svc.Repo.Create(&newTrack); err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed to save track")
			return
		}

		newTrackID = newTrack.ID
		log.WithFields(log.Fields{"track_id": newTrackID}).Info("Track saved successfully")
	}()

	wg.Wait()

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

	// update physical ptr of file defined by uuid
	if uuid != "" {
		if err := ctl.Svc.Repo.UpdateTrackPtr(newTrackID, uuid); err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed to update track UUID ptr")
		} else {
			log.WithFields(log.Fields{"uuid": uuid}).Info("Track UUID ptr updated successfully")
		}
	}
}

func (ctl *MusicCtl) RemoveTrack(c *gin.Context) {
	trackID, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

	if err := ctl.Svc.Repo.DeleteByID(trackID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete track"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Track deleted successfully"})
}

func (ctl *MusicCtl) UpdateTrackMetadata(c *gin.Context) {
	trackID, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }
	
	var req UpdateTrackMetadataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctl.Svc.UpdateTrackMetadata(trackID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update track metadata"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Track metadata updated successfully"})
}

// TODO pagination
func (ctl *MusicCtl) ListTracks(c *gin.Context) {
	tracks, err := ctl.Svc.Repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tracks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tracks": tracks})
}

func (ctl *MusicCtl) GetTrackByID(c *gin.Context) {
    trackID, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    track, err := ctl.Svc.Repo.GetByID(trackID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Track not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"track": track})
}

func (ctl *MusicCtl) PlayTrack(c *gin.Context) {
    // Получаем трек по ID или UUID
    trackID, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
        return
    }

    track, err := ctl.Svc.Repo.GetByID(trackID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Track not found"})
        return
    }

    if track.Ptr == "" {
        c.JSON(http.StatusNotFound, gin.H{"error": "Track file not found"})
        return
    }

    // Возвращаем ссылку на трек по UUID
    trackURL := "/upload/" + track.Ptr
    c.JSON(http.StatusOK, gin.H{"track_url": trackURL})
}
