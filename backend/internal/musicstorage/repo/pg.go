package repo

import (
	. "untitled/internal/musicstorage/mdl"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PgRepo struct {
	db *gorm.DB
}

func NewPgRepo(db *gorm.DB) *PgRepo {
	return &PgRepo{db: db}
}

func (r *PgRepo) CreateTrack(track *Track_PG) error {
	return r.db.Create(track).Error
}

func (r *PgRepo) DeleteTrackByID(trackID uuid.UUID) error {
	return r.db.Where("uuid = ?", trackID).Delete(&Track_PG{}).Error
}

func (r *PgRepo) UpdateTrackMetadata(trackID uuid.UUID) error {
	return r.db.Model(&Track_PG{}).Where("uuid = ?", trackID).Update("updated_at", gorm.Expr("NOW()")).Error
}

func (r *PgRepo) GetTrackByID(trackID uuid.UUID) (*Track_PG, error) {
	var track Track_PG
	err := r.db.Where("uuid = ?", trackID).First(&track).Error
	if err != nil {
		return nil, err
	}
	return &track, nil
}

func (r *PgRepo) ListTracks() ([]Track_PG, error) {
	var tracks []Track_PG
	err := r.db.Find(&tracks).Error
	return tracks, err
}
