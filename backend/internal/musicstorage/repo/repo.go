package repo

import (
	. "untitled/internal/musicstorage/mdl"

	"gorm.io/gorm"
)

type TrackRepo struct {
	db *gorm.DB
}

func NewTrackRepo(db *gorm.DB) *TrackRepo {
	return &TrackRepo{db: db}
}

func (r *TrackRepo) Create(track *Track) error {
	return r.db.Create(track).Error
}

func (r *TrackRepo) GetByID(id int64) (*Track, error) {
	var track Track
	err := r.db.First(&track, id).Error
	return &track, err
}

func (r *TrackRepo) GetAll() ([]Track, error) {
	var tracks []Track
	err := r.db.Find(&tracks).Error
	return tracks, err
}

func (r *TrackRepo) Update(track *Track) error {
	return r.db.Save(track).Error
}

func (r *TrackRepo) DeleteByID(id int64) error {
	return r.db.Delete(&Track{}, id).Error
}

func (r *TrackRepo) UpdateTrackPtr(id int64, ptr string) error {
	return r.db.Model(&Track{}).Where("id = ?", id).Update("ptr", ptr).Error
}
