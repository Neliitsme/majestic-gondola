package models

import "time"

type Artist struct {
	Id        int `pg:"artist_id,pk"`
	Name      string
	CreatedAt time.Time `pg:"default:now()"`
	Score     int       `pg:"default:0"`
}
