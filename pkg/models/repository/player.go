package repository

import "time"

type Player struct {
	Puuid     string    `db:"puuid"`
	Username  string    `db:"username"`
	Tag       string    `db:"tag"`
	CreatedAt time.Time `db:"created_at"`
}
