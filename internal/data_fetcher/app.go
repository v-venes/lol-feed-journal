package datafetcher

import (
	"context"
	"encoding/json"
	"log"
	"slices"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	redismodel "github.com/v-venes/feed_journal/pkg/models/redis"
	"github.com/v-venes/feed_journal/pkg/models/repository"
	"github.com/v-venes/feed_journal/pkg/repositories"
	"github.com/v-venes/feed_journal/pkg/services"
)

type Application struct {
	LeagueService    *services.LeagueService
	PlayerRepository *repositories.PlayerRepository
	MatchRepository  *repositories.MatchRepository
	RedisClient      *redis.Client
	RedisChannel     string
}

type NewApplicationParams struct {
	RiotBasePath     string
	RiotDDBasePath   string
	RiotApiKey       string
	RedisChannel     string
	PlayerRepository *repositories.PlayerRepository
	MatchRepository  *repositories.MatchRepository
	RedisClient      *redis.Client
}

func NewApplication(params NewApplicationParams) *Application {
	leagueService := services.NewLeagueService(services.NewLeagueServiceParams{
		Key:        params.RiotApiKey,
		BasePath:   params.RiotBasePath,
		DDBasePath: params.RiotDDBasePath,
	})

	return &Application{
		LeagueService:    leagueService,
		PlayerRepository: params.PlayerRepository,
		MatchRepository:  params.MatchRepository,
		RedisClient:      params.RedisClient,
		RedisChannel:     params.RedisChannel,
	}
}

func (a *Application) Run() {
	brLoc, _ := time.LoadLocation("America/Sao_Paulo")
	fetchDate := time.Now().In(brLoc)
	fetchStartDate := time.Date(fetchDate.Year(), fetchDate.Month(), fetchDate.Day(), 0, 0, 0, 0, fetchDate.Location())
	fetchEndDate := time.Date(fetchDate.Year(), fetchDate.Month(), fetchDate.Day(), 23, 59, 59, 0, fetchDate.Location())

	players, err := a.PlayerRepository.GetAll()

	if err != nil {
		log.Fatal(err)
	}

	a.startFetchData(players, fetchStartDate, fetchEndDate)

	a.sendToGenerateImageQueue(fetchStartDate)
}

func (a *Application) startFetchData(players []repository.Player, fetchStartDate time.Time, fetchEndDate time.Time) {

	for _, player := range players {
		log.Printf("Fetching match data for %s\n", player.Username)

		matchesIds, err := a.LeagueService.GetMatchesIDs(
			services.GetMatchesIDsParams{AccountID: player.Puuid, From: uint32(fetchStartDate.UnixMilli() / 1000), To: uint32(fetchEndDate.UnixMilli() / 1000)},
		)

		if err != nil {
			log.Fatal(err)
		}

		var wg sync.WaitGroup
		for _, matchID := range matchesIds {
			wg.Add(1)
			go func(matchID string) {
				defer wg.Done()
				err := a.processMatch(matchID, fetchStartDate)
				if err != nil {
					log.Fatal(err)
				}
			}(matchID)
		}

		wg.Wait()

		log.Printf("Fetch completed for %s\n", player.Username)

		if len(matchesIds) > 0 {
			time.Sleep(30 * time.Second)
		}
	}

}

func (a *Application) processMatch(matchID string, fetchDate time.Time) error {
	matchDetails, err := a.LeagueService.GetMatchDetails(matchID)
	gameModesToIgnore := []string{"URF", "SWIFTPLAY"}

	if err != nil {
		return err
	}

	if slices.Contains(gameModesToIgnore, matchDetails.Info.GameMode) {
		return nil
	}

	matchsByPlayers := matchDetails.ToRepositoryMatch(fetchDate)

	if len(matchsByPlayers) == 0 {
		return nil
	}

	err = a.MatchRepository.SaveMatchs(matchsByPlayers)

	if err != nil {
		return err
	}

	return nil
}

func (a *Application) sendToGenerateImageQueue(matchDate time.Time) error {
	ctx := context.Background()

	message := redismodel.GenerateImagePayload{MatchDate: matchDate}
	jsonMsg, _ := json.Marshal(message)

	err := a.RedisClient.Publish(ctx, a.RedisChannel, jsonMsg).Err()

	if err != nil {
		return err
	}

	return nil

}
