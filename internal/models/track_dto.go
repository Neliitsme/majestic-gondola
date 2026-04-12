package models

import "time"

type CreateTrackRequest struct {
	Name        string    `json:"name" binding:"required"`
	Author      string    `json:"author" binding:"required"`
	ReleaseDate time.Time `json:"release_date"`
	Genres      []string  `json:"genres"`
}
