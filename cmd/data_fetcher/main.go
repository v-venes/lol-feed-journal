package main

import (
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	datafetcher "github.com/v-venes/feed_journal/internal/data_fetcher"
	"github.com/v-venes/feed_journal/pkg/config"
	"github.com/v-venes/feed_journal/pkg/repositories"
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

	playerRepository := repositories.NewPlayerRepository(db)
	matchRepository := repositories.NewMatchRepository(db)

	app := datafetcher.NewApplication(datafetcher.NewApplicationParams{
		RiotBasePath:     env.RiotBasePath,
		RiotApiKey:       env.RiotKey,
		PlayerRepository: playerRepository,
		MatchRepository:  matchRepository,
		RedisClient:      rdb,
		RedisChannel:     env.RedisChannel,
	})

	app.Run()
}
