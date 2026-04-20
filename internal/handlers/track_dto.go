package handlers

type CreateTrackRequest struct {
	Name        string   `json:"name" binding:"required" example:"Bohemian Rhapsody"`
	Author      string   `json:"author" binding:"required" example:"Queen"`
	ReleaseDate *string  `json:"release_date" example:"1975-10-31 15:04:05"`
	Genres      []string `json:"genres" example:"Rock,Opera"`
}

type UpdateTrackRequest struct {
	Id          int      `json:"id" binding:"required" example:"1"`
	Name        string   `json:"name" binding:"required" example:"Bohemian Rhapsody"`
	Author      string   `json:"author" binding:"required" example:"Queen"`
	ReleaseDate *string  `json:"release_date" example:"1975-10-31 15:04:05"`
	Genres      []string `json:"genres" example:"Rock,Opera"`
}

type TrackResponse struct {
	Id          int      `json:"id" example:"1"`
	Name        string   `json:"name" example:"Bohemian Rhapsody"`
	Author      string   `json:"author" example:"Queen"`
	ReleaseDate string   `json:"release_date" example:"1975-10-31 15:04:05"`
	Genres      []string `json:"genres" example:"Rock,Opera"`
	CreatedAt   string   `json:"created_at" example:"2006-01-02 15:04:05"`
}
