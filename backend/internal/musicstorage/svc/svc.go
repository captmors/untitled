package svc

import (
	"errors"
	"strings"
	. "untitled/internal/musicstorage/repo"
	. "untitled/internal/musicstorage/mdl"

	"github.com/gin-gonic/gin"
)

type MusicSvc struct {
	Repo *TrackRepo
}

func NewMusicSvc(repo *TrackRepo) *MusicSvc {
	return &MusicSvc{Repo: repo}
}

func GetUUIDFromLocation(location string) string {
	parts := strings.Split(location, "/")
	return parts[len(parts)-1]
}

func GetCurrentUserID(c *gin.Context) (int64, error) {
	userID, ok := c.Get("UserID")
	if !ok {
		return 0, errors.New("UserID not found")
	}

	id, ok := userID.(int64)
	if !ok {
		return 0, errors.New("invalid UserID type")
	}

	return id, nil
}

func (svc *MusicSvc) UpdateTrackMetadata(trackID int64, req UpdateTrackMetadataRequest) error {
	track, err := svc.Repo.GetByID(trackID)
	if err != nil {
		return err
	}

	if req.Title != "" {
		track.Title = req.Title
	}
	if req.Artist != "" {
		track.Artist = req.Artist
	}
	if req.Format != "" {
		track.Format = req.Format
	}
	if req.Bitrate != nil {
		track.Bitrate = req.Bitrate
	}
	if req.Duration != nil {
		track.Duration = req.Duration
	}
	if req.Genre != nil {
		track.Genre = req.Genre
	}

	return svc.Repo.Update(track)
}
