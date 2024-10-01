package mdl

type TrackRequest struct {
	Title    string  `json:"title,omitempty" binding:"required"`
	Artist   string  `json:"artist,omitempty" binding:"required"`
	Format   string  `json:"format,omitempty" binding:"required"`
	Duration *int    `json:"duration,omitempty"`
	Album    *string `json:"album,omitempty"`
	Genre    *string `json:"genre,omitempty"`
	Ptr      *string `json:"ptr,omitempty"`
}

type ListTracksResponse struct {
	Tracks []TrackResponse `json:"tracks"`
}

type TrackResponse struct {
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Duration *int   `json:"duration,omitempty"`
}
