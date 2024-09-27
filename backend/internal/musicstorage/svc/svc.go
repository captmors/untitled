package svc

import (
	"errors"
	"strings"
	. "untitled/internal/musicstorage/repo"

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

func GetCurrentUserID(c *gin.Context) (uint, error) {
	userID, ok := c.Get("UserID")
	if !ok {
		return 0, errors.New("UserID not found")
	}

	id, ok := userID.(uint)
	if !ok {
		return 0, errors.New("invalid UserID type")
	}

	return id, nil
}
