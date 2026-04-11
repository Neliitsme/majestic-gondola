package repository

import (
	"time"
)

type Song struct {
	Id          int `pg:"pk_id"`
	Name        string
	Author      string
	ReleaseDate time.Time
	Genres      []string
	CreatedAt   time.Time `pg:"default:now()"`
}
