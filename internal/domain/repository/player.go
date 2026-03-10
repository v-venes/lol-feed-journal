package domain

import "time"

type Player struct {
	Puuid        string        `gorm:"primaryKey;type:varchar(128)"`
	Username     string        `gorm:"type:varchar(16);not null"`
	Tag          string        `gorm:"type:varchar(16);not null"`
	CreatedAt    time.Time     `gorm:"typetimestamptz;not null"`
	MatchPlayers []MatchPlayer `gorm:"foreignKey:Puuid"`
}

func (Player) TableName() string {
	return "players"
}
