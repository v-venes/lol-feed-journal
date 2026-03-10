package domain

import "time"

type Match struct {
	Id           string        `gorm:"primaryKey;type:varchar(16)"`
	GameMode     string        `gorm:"type:varchar(32);not null"`
	GameType     string        `gorm:"type:varchar(32);not null"`
	MatchCreated time.Time     `gorm:"typetimestamptz;not null"`
	GameStarted  time.Time     `gorm:"typetimestamptz;not null"`
	GameEnded    time.Time     `gorm:"typetimestamptz;not null"`
	MatchPlayers []MatchPlayer `gorm:"foreignKey:MatchID"`
}

func (Match) TableName() string {
	return "matches"
}
