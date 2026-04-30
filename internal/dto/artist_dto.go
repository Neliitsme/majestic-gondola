package dto

type CreateArtistRequest struct {
	Name string `json:"name" binding:"required" example:"Queen"`
}

type UpdateArtistRequest struct {
	Name string `json:"name" binding:"required" example:"Queen"`
}

type ArtistResponse struct {
	Id        int    `json:"id" example:"1"`
	Name      string `json:"name" example:"Queen"`
	CreatedAt string `json:"created_at" example:"2006-01-02 15:04:05"`
	Score     int    `json:"score" example:"100"`
}
