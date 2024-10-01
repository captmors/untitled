package svc

import (
	"errors"
	"strings"
	"time"
	. "untitled/internal/musicstorage/mdl"
	. "untitled/internal/musicstorage/repo"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MusicSvc struct {
	PgRepo    *PgRepo    
	MongoRepo *MongoRepo 
}

func NewMusicSvc(pgRepo *PgRepo, mongoRepo *MongoRepo) *MusicSvc {
	return &MusicSvc{
		PgRepo:    pgRepo,
		MongoRepo: mongoRepo,
	}
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

func (svc *MusicSvc) CreateTrack(req TrackRequest, userID int64) error {
	now := time.Now()

	newTrackPg := Track_PG{
		UUID:      uuid.New(),
		UserID:    userID,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	newTrackMongo := Track_MONGO{
		ID:       newTrackPg.UUID, 
		Title:    req.Title,
		Artist:   req.Artist,
		Album:    req.Album,
		Format:   req.Format,
		Duration: req.Duration,
		Genre:    req.Genre,
	}

	if err := svc.PgRepo.CreateTrack(&newTrackPg); err != nil {
		return err
	}
	if err := svc.MongoRepo.CreateTrack(newTrackMongo); err != nil {
		return err
	}

	return nil
}

func (svc *MusicSvc) UpdateTrackPtrByUUID(trackUUID string) error {
	if err := svc.MongoRepo.UpdateTrackPtr(trackUUID); err != nil {
		return err
	}
	return nil
}

func (svc *MusicSvc) DeleteTrackByID(trackID uuid.UUID) error {
	if err := svc.PgRepo.DeleteTrackByID(trackID); err != nil {
		return err
	}
	if err := svc.MongoRepo.DeleteTrackByID(trackID); err != nil {
		return err
	}
	return nil
}

func (svc *MusicSvc) UpdateTrackMetadata(trackID uuid.UUID, req TrackRequest) error {
	if err := svc.PgRepo.UpdateTrackMetadata(trackID); err != nil {
		return err
	}

	if err := svc.MongoRepo.UpdateTrackMetadata(trackID, req); err != nil {
		return err
	}

	return nil
}

func (svc *MusicSvc) ListTracks() ([]*Track, error) {
	tracksPG, err := svc.PgRepo.ListTracks()
	if err != nil {
		return nil, err
	}

	var combinedTracks []*Track

	for _, trackPG := range tracksPG {
		trackMongo, err := svc.MongoRepo.GetTrackByID(trackPG.UUID)
		if err != nil {
			continue
		}
		combinedTrack := &Track{
			Track_PG:   trackPG,
			Track_MONGO: *trackMongo,
		}
		combinedTracks = append(combinedTracks, combinedTrack)
	}

	return combinedTracks, nil
}

func (svc *MusicSvc) GetTrackByID(trackID uuid.UUID) (*Track, error) {
	trackPG, err := svc.PgRepo.GetTrackByID(trackID)
	if err != nil {
		return nil, err
	}

	trackMongo, err := svc.MongoRepo.GetTrackByID(trackID)
	if err != nil {
		return nil, err
	}

	return &Track{
		Track_PG:   *trackPG,
		Track_MONGO: *trackMongo,
	}, nil
}
