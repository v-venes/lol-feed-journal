package feedjournal

import (
	"time"

	"github.com/v-venes/feed_journal/pkg/models/repository"
)

type Journal struct {
	StoredPlayers []repository.Match
	RandomPlayers []repository.Match
	MatchDate     time.Time
}
