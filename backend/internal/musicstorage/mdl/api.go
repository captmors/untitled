package mdl

type UploadTrackRequest struct {
    Title    string  `json:"title" binding:"required"`
    Artist   string  `json:"artist" binding:"required"`
    Format   string  `json:"format" binding:"required"`
    Bitrate  *int    `json:"bitrate,omitempty"`
    Duration *int    `json:"duration,omitempty"`
    Genre    *string `json:"genre,omitempty"`
}

type UpdateTrackMetadataRequest struct {
	Title   string  `json:"title,omitempty"`
	Artist  string  `json:"artist,omitempty"`
	Format  string  `json:"format,omitempty"`
	Bitrate *int    `json:"bitrate,omitempty"`
	Duration *int   `json:"duration,omitempty"`
	Genre   *string `json:"genre,omitempty"`
}
