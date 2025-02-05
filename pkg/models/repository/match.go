package repository

import "time"

type Match struct {
	Id                 string    `db:"id"`
	Puuid              string    `db:"puuid"`
	ProfileIcon        uint32    `db:"profile_icon"`
	SummonerName       string    `db:"summoner_name"`
	Assists            uint8     `db:"assists"`
	Deaths             uint8     `db:"deaths"`
	Kills              uint8     `db:"kills"`
	KDA                float32   `db:"kda"`
	ChampionName       string    `db:"champion_name"`
	Lane               string    `db:"lane"`
	Role               string    `db:"role"`
	IndividualPosition string    `db:"individual_position"`
	TeamPosition       string    `db:"team_position"`
	Item0              uint32    `db:"item0"`
	Item1              uint32    `db:"item1"`
	Item2              uint32    `db:"item2"`
	Item3              uint32    `db:"item3"`
	Item4              uint32    `db:"item4"`
	Item5              uint32    `db:"item5"`
	Item6              uint32    `db:"item6"`
	Win                bool      `db:"win"`
	Surrender          bool      `db:"surrender"`
	Remake             bool      `db:"remake"`
	MatchDate          time.Time `db:"match_date"`
}
