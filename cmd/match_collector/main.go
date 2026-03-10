package main

import (
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/v-venes/lol-feed-journal/internal/config"
	matchcollector "github.com/v-venes/lol-feed-journal/internal/match_collector"
	"github.com/v-venes/lol-feed-journal/internal/repository"
)

func init() {
	godotenv.Load()
}

func main() {
	env := config.GetConfigVars()

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		env.PostgresUser,
		env.PostgresPassword,
		env.PostgresHost,
		env.PostgresDB,
	)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: env.RedisHost,
		DB:   0,
	})
	defer rdb.Close()

	playerRepository := repository.NewPlayerRepository(db)
	matchRepository := repository.NewMatchRepository(db)

	app := matchcollector.NewApplication(matchcollector.NewApplicationParams{
		RiotBasePath:     env.RiotBasePath,
		RiotApiKey:       env.RiotKey,
		RiotDDBasePath:   env.RiotDDBasePath,
		PlayerRepository: playerRepository,
		MatchRepository:  matchRepository,
		RedisClient:      rdb,
		RedisChannel:     env.RedisChannel,
	})

	app.Run()
}
