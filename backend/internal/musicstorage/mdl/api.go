package mdl

type UploadTrackRequest struct {
	Title  string `form:"title" binding:"required"`
	Artist string `form:"artist" binding:"required"`
	Format string `form:"format" binding:"required"` 
}