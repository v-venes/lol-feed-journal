package league

import (
	"time"

	"github.com/v-venes/lol-feed-journal/pkg/models/repository"
)

type MatchChallenges struct {
	KDA float32 `json:"kda"`
}

type MatchParticipants struct {
	Puuid              string          `json:"puuid"`
	Assists            uint8           `json:"assists"`
	Deaths             uint8           `json:"deaths"`
	Kills              uint8           `json:"kills"`
	Challenges         MatchChallenges `json:"challenges"`
	ChampionName       string          `json:"championName"`
	Lane               string          `json:"lane"`
	Role               string          `json:"role"`
	IndividualPosition string          `json:"individualPosition"`
	TeamPosition       string          `json:"teamPosition"`
	Item0              uint32          `json:"item0"`
	Item1              uint32          `json:"item1"`
	Item2              uint32          `json:"item2"`
	Item3              uint32          `json:"item3"`
	Item4              uint32          `json:"item4"`
	Item5              uint32          `json:"item5"`
	Item6              uint32          `json:"item6"`
	ProfileIcon        uint32          `json:"profileIcon"`
	SummonerName       string          `json:"riotIdGameName"`
	Win                bool            `json:"win"`
	Surrender          bool            `json:"gameEndedInSurrender"`
	Remake             bool            `json:"gameEndedInEarlySurrender"`
}

type MatchMetadata struct {
	MatchID string `json:"matchId"`
}

type MatchInfo struct {
	GameMode     string              `json:"gameMode"`
	Participants []MatchParticipants `json:"participants"`
}

type Match struct {
	Metadata MatchMetadata `json:"metadata"`
	Info     MatchInfo     `json:"info"`
}

func (m *Match) ToRepositoryMatch(matchDate time.Time) []repository.Match {
	if m == nil {
		return nil
	}

	var matchs []repository.Match

	for _, participant := range m.Info.Participants {

		if participant.Remake {
			continue
		}

		kills := float64(participant.Kills)
		assists := float64(participant.Assists)
		deaths := float64(participant.Deaths)

		kda := kills + (deaths * -2) + (assists * 0.5)
		if participant.Role == "SUPPORT" && participant.TeamPosition == "UTILITY" {
			kda = (kills * 0.75) + (deaths * -2) + (assists * 0.75)
		}

		matchs = append(matchs, repository.Match{
			Id:                 m.Metadata.MatchID,
			Puuid:              participant.Puuid,
			ProfileIcon:        participant.ProfileIcon,
			SummonerName:       participant.SummonerName,
			Assists:            participant.Assists,
			Deaths:             participant.Deaths,
			Kills:              participant.Kills,
			KDA:                float32(kda),
			ChampionName:       participant.ChampionName,
			Lane:               participant.Lane,
			Role:               participant.Role,
			IndividualPosition: participant.IndividualPosition,
			TeamPosition:       participant.TeamPosition,
			Item0:              participant.Item0,
			Item1:              participant.Item1,
			Item2:              participant.Item2,
			Item3:              participant.Item3,
			Item4:              participant.Item4,
			Item5:              participant.Item5,
			Item6:              participant.Item6,
			Win:                participant.Win,
			Surrender:          participant.Surrender,
			Remake:             participant.Remake,
			MatchDate:          matchDate,
		})

	}

	return matchs
}
