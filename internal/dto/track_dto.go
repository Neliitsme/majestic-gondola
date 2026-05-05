package dto

type CreateTrackRequest struct {
	Name        string   `json:"name" binding:"required" example:"Bohemian Rhapsody"`
	ArtistId    *int     `json:"artist_id" example:"1"`
	ReleaseDate *string  `json:"release_date" example:"1975-10-31 15:04:05"`
	Genres      []string `json:"genres" example:"Rock,Opera"`
}

type UpdateTrackRequest struct {
	Name        string   `json:"name" binding:"required" example:"Bohemian Rhapsody"`
	ArtistId    *int     `json:"artist_id" example:"1"`
	ReleaseDate *string  `json:"release_date" example:"1975-10-31 15:04:05"`
	Genres      []string `json:"genres" example:"Rock,Opera"`
}

type TrackResponse struct {
	Id          int      `json:"id" example:"1"`
	Name        string   `json:"name" example:"Bohemian Rhapsody"`
	ArtistId    *int     `json:"artist_id" example:"1"`
	ReleaseDate string   `json:"release_date" example:"1975-10-31 15:04:05"`
	Genres      []string `json:"genres" example:"Rock,Opera"`
	CreatedAt   string   `json:"created_at" example:"2006-01-02 15:04:05"`
	Score       int      `json:"score" example:"100"`
	ReviewCount int      `json:"review_count" example:"100"`
}

type PopulateRequest struct {
	ArtistId *int `json:"artist_id" example:"1"`
}
