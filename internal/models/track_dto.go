package models

import "time"

type CreateTrackRequest struct {
	Name        string    `json:"name" binding:"required"`
	Author      string    `json:"author" binding:"required"`
	ReleaseDate time.Time `json:"release_date"`
	Genres      []string  `json:"genres"`
}

type UpdateTrackRequest struct {
	Id          int       `json:"id" binding:"required"`
	Name        string    `json:"name"`
	Author      string    `json:"author"`
	ReleaseDate time.Time `json:"release_date"`
	Genres      []string  `json:"genres"`
}
