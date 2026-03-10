package feedjournal

import (
	"time"

	domain "github.com/v-venes/lol-feed-journal/internal/domain/repository"
)

type Journal struct {
	StoredPlayers []domain.Match
	RandomPlayers []domain.Match
	MatchDate     time.Time
}
