package mdl

import (
	"time"

	"github.com/google/uuid"
)

type Track struct {
	Track_PG
	Track_MONGO
}

type Track_PG struct {
	UUID   uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"uuid"` // Common UUID
	UserID int64     `json:"user_id"`                                                      // UserID from Track_PG

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type Track_MONGO struct {
	ID       uuid.UUID `bson:"_id" json:"id"`
	Title    string    `bson:"title" json:"title"`
	Artist   string    `bson:"artist" json:"artist"`
	Album    *string   `bson:"album,omitempty" json:"album,omitempty"`
	Genre    *string   `bson:"genre,omitempty" json:"genre,omitempty"`
	Format   string    `bson:"format" json:"format"`
	Duration *int      `bson:"duration,omitempty" json:"duration,omitempty"`
	Ptr      *string   `bson:"ptr,omitempty" json:"ptr,omitempty"`
}
