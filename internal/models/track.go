package models

import (
	"time"
)

type Track struct {
	Id          int `pg:",pk"`
	Name        string
	Author      string
	ReleaseDate time.Time `pg:"default:now()"`
	Genres      []string  `pg:",array"`
	CreatedAt   time.Time `pg:"default:now()"`
}
