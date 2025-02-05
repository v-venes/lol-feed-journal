package repositories

import (
	"time"

	"github.com/jmoiron/sqlx"
	repositorymodel "github.com/v-venes/feed_journal/pkg/models/repository"
)

type MatchRepository struct {
	*Repository[repositorymodel.Match]
}

func NewMatchRepository(db *sqlx.DB) *MatchRepository {
	return &MatchRepository{
		Repository: NewRepository[repositorymodel.Match](db),
	}
}

func (p *MatchRepository) SaveMatchs(matchs []repositorymodel.Match) error {
	query := `
        INSERT INTO players_matches (
            id, puuid, profile_icon, summoner_name, assists, deaths, kills, kda,
            champion_name, lane, role, individual_position, team_position,
            item0, item1, item2, item3, item4, item5, item6,
            win, surrender, remake, match_date
        )
        VALUES (
            :id, :puuid, :profile_icon, :summoner_name, :assists, :deaths, :kills, :kda,
            :champion_name, :lane, :role, :individual_position, :team_position,
            :item0, :item1, :item2, :item3, :item4, :item5, :item6,
            :win, :surrender, :remake, :match_date
        )
        ON CONFLICT (id, puuid) DO NOTHING
    `
	err := p.Repository.SaveMany(query, matchs)

	if err != nil {
		return err
	}

	return nil
}

func (p *MatchRepository) GetTop10StoredPlayersMatches(matchDate time.Time) ([]repositorymodel.Match, error) {
	query := `
        SELECT * FROM players_matches  
        WHERE puuid IN (SELECT puuid FROM players) 
        AND match_date = $1::timestamp  
        ORDER BY kda ASC 
        LIMIT 10
    `
	matches, err := p.Repository.GetAll(query, matchDate)

	if err != nil {
		return nil, err
	}

	return matches, nil
}

func (p *MatchRepository) GetTop10RandomPlayersMatches(matchDate time.Time) ([]repositorymodel.Match, error) {
	query := `
        SELECT * FROM players_matches  
        WHERE puuid NOT IN (SELECT puuid FROM players) 
        AND match_date = $1::timestamp  
        ORDER BY kda ASC 
        LIMIT 10
    `
	matches, err := p.Repository.GetAll(query, matchDate)

	if err != nil {
		return nil, err
	}

	return matches, nil
}
