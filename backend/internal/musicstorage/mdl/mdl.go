package mdl

import "gorm.io/gorm"

type Track struct {
	gorm.Model
	ID       int64   `gorm:"primaryKey" json:"id"`
	Title    string  `gorm:"not null" json:"title"`
	Artist   string  `gorm:"not null" json:"artist"`
	Format   string  `gorm:"not null" json:"format"`
	Ptr      string  `gorm:"default:null" json:"ptr,omitempty"`
	Bitrate  *int    `gorm:"default:null" json:"bitrate,omitempty"`
	Duration *int    `gorm:"default:null" json:"duration,omitempty"`
	Genre    *string `gorm:"default:null" json:"genre,omitempty"`
	UserID   int64   `json:"user_id"`
}
