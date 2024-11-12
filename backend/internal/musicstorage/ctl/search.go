package ctl

import (
	. "untitled/internal/musicstorage/mdl"

	"github.com/gin-gonic/gin"
)

func (ctl *MusicCtl) SearchTracks(c *gin.Context) {
	var req TrackSearchRequest
	if err := c.BindQuery(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tracks, err := ctl.Svc.ESRepo.SearchTracks(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, tracks)
}
