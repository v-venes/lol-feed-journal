package domain

type MatchPlayer struct {
	MatchID            string `gorm:"primaryKey;type:varchar(16)"`
	Puuid              string `gorm:"primaryKey;type:varchar(128)"`
	ChampionID         string `gorm:"type:uint;not null"`
	ChampionName       string `gorm:"type:varchar(64);not null"`
	ChampionLevel      uint8  `gorm:"type:uint;not null"`
	Kills              uint8  `gorm:"type:uint;not null"`
	Deaths             uint8  `gorm:"type:uint;not null"`
	Assists            uint8  `gorm:"type:uint;not null"`
	Lane               string `gorm:"type:varchar(16)"`
	IndividualPosition string `gorm:"type:varchar(16)"`
	Role               string `gorm:"type:varchar(16)"`
	TeamPosition       string `gorm:"type:varchar(16)"`
	Spell1ID           uint8  `gorm:"type:uint;not null"`
	Spell2ID           uint8  `gorm:"type:uint;not null"`
	Item0              uint8  `gorm:"type:uint;not null"`
	Item1              uint8  `gorm:"type:uint;not null"`
	Item2              uint8  `gorm:"type:uint;not null"`
	Item3              uint8  `gorm:"type:uint;not null"`
	Item4              uint8  `gorm:"type:uint;not null"`
	Item5              uint8  `gorm:"type:uint;not null"`
	Item6              uint8  `gorm:"type:uint;not null"`
	Win                bool   `gorm:"type:bool;not null"`
	Player             Player `gorm:"foreignKey:Puuid"`
	Match              Match  `gorm:"foreignKey:MatchID"`
}

func (MatchPlayer) TableName() string {
	return "matches_players"
}
