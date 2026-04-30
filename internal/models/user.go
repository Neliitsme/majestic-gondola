package models

import "time"

type User struct {
	Id        int `pg:"user_id,pk"`
	Name      string
	CreatedAt time.Time `pg:"default:now()"`
}
