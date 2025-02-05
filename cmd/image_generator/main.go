package main

import (
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	imagegenerator "github.com/v-venes/feed_journal/internal/image_generator"
	"github.com/v-venes/feed_journal/pkg/config"
	"github.com/v-venes/feed_journal/pkg/repositories"
	"github.com/v-venes/feed_journal/pkg/services"
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

	minioClient, err := minio.New(env.MinioHost, &minio.Options{
		Creds:  credentials.NewStaticV4(env.MinioKey, env.MinioSecret, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal(err)
	}

	leagueService := services.NewLeagueService(services.NewLeagueServiceParams{
		Key:      env.RiotKey,
		BasePath: env.RiotBasePath,
	})

	playerRepository := repositories.NewPlayerRepository(db)
	matchRepository := repositories.NewMatchRepository(db)

	app := imagegenerator.NewApplication(imagegenerator.NewApplicationParams{
		PlayerRepository: playerRepository,
		MatchRepository:  matchRepository,
		LeagueService:    leagueService,
		RedisClient:      rdb,
		RedisChannel:     env.RedisChannel,
		MinioClient:      minioClient,
	})

	app.Run()
}
