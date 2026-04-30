package models

import "time"

type Review struct {
	Id          int `pg:"review_id,pk"`
	UserId      *int
	TrackId     *int
	Score       int       `pg:"default:0"`
	IsDeleted   bool      `pg:"default:false"`
	IsProcessed bool      `pg:"default:false"`
	CreatedAt   time.Time `pg:"default:now()"`
	User        *User     `pg:"rel:has-one"`
	Track       *Track    `pg:"rel:has-one"`
}
