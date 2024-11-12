package ctl

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/tus/tusd/v2/pkg/handler"

	. "untitled/internal/musicstorage/mdl"
	. "untitled/internal/musicstorage/svc"

	"github.com/google/uuid"
)

type MusicCtl struct {
	Svc        *MusicSvc
	tusHandler *handler.UnroutedHandler
}

func NewMusicCtl(musicSvc *MusicSvc, tusHandler *handler.UnroutedHandler) *MusicCtl {
	return &MusicCtl{Svc: musicSvc, tusHandler: tusHandler}
}

func (ctl *MusicCtl) UploadTrack(c *gin.Context) {
	var req TrackRequest
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
	var trackUUID string

	// upload file resumably
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("Starting file upload")
		ctl.tusHandler.PostFile(c.Writer, c.Request)
		uploadStatus = c.Writer.Status()
		trackUUID = GetUUIDFromLocation(c.Writer.Header().Get("Location"))
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

		if err := ctl.Svc.CreateTrack(req, userID); err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed to save track")
			return
		}
		log.Info("Track saved successfully")
	}()

	wg.Wait()

	if uploadStatus != http.StatusCreated {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed"})
		log.Error("Upload failed with status:", uploadStatus)
		return
	}

	// update physical ptr of file defined by uuid
	if trackUUID != "" {
		if err := ctl.Svc.UpdateTrackPtrByUUID(trackUUID); err != nil {
			log.WithFields(log.Fields{"error": err}).Error("Failed to update track UUID ptr")
		} else {
			log.WithFields(log.Fields{"uuid": trackUUID}).Info("Track UUID ptr updated successfully")
		}
	}
}

func (ctl *MusicCtl) RemoveTrack(c *gin.Context) {
	trackIDStr := c.Param("id")
	trackID, err := uuid.Parse(trackIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	if err := ctl.Svc.DeleteTrackByID(trackID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete track"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Track deleted successfully"})
}

func (ctl *MusicCtl) UpdateTrackMetadata(c *gin.Context) {
	trackIDStr := c.Param("id")
	trackID, err := uuid.Parse(trackIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var req TrackRequest
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
	tracks, err := ctl.Svc.ListTracks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tracks"})
		return
	}

	var responseTracks []TrackResponse
	for _, track := range tracks {
		responseTracks = append(responseTracks, TrackResponse{
			Title:    track.Title,
			Artist:   track.Artist,
			Duration: track.Duration,
		})
	}

	c.JSON(http.StatusOK, ListTracksResponse{Tracks: responseTracks})
}

func (ctl *MusicCtl) GetTrackByID(c *gin.Context) {
	trackIDStr := c.Param("id")
	trackID, err := uuid.Parse(trackIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	track, err := ctl.Svc.GetTrackByID(trackID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Track not found"})
		return
	}

	c.JSON(http.StatusOK, TrackResponse{
		Title:    track.Title,
		Artist:   track.Artist,
		Duration: track.Duration,
	})
}

func (ctl *MusicCtl) PlayTrack(c *gin.Context) {
	trackIDStr := c.Param("id")
	trackID, err := uuid.Parse(trackIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	track, err := ctl.Svc.GetTrackByID(trackID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Track not found"})
		return
	}

	if track.Ptr == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Track file not found"})
		return
	}

	trackURL := "/upload/" + *track.Ptr
	c.JSON(http.StatusOK, gin.H{"track_url": trackURL})
}
