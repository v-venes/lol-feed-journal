package repositories

import (
	"github.com/jmoiron/sqlx"
	repositorymodel "github.com/v-venes/feed_journal/pkg/models/repository"
)

type PlayerRepository struct {
	*Repository[repositorymodel.Player]
}

func NewPlayerRepository(db *sqlx.DB) *PlayerRepository {
	return &PlayerRepository{
		Repository: NewRepository[repositorymodel.Player](db),
	}
}

func (p *PlayerRepository) GetAll() ([]repositorymodel.Player, error) {
	players, err := p.Repository.GetAll("SELECT * FROM players")

	if err != nil {
		return nil, err
	}

	return players, nil
}
