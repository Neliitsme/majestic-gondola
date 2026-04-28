package models

import (
	"time"
)

type Track struct {
	Id          int `pg:"track_id,pk"`
	Name        string
	ArtistId    *int
	ReleaseDate time.Time `pg:"default:now()"`
	Genres      []string  `pg:",array"`
	CreatedAt   time.Time `pg:"default:now()"`
	Score       int       `pg:"default:0"`
	Artist      *Artist   `pg:"rel:has-one"`
}
