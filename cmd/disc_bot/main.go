package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	discordbot "github.com/v-venes/feed_journal/internal/disc_bot"
	"github.com/v-venes/feed_journal/pkg/config"
)

func init() {
	godotenv.Load()
}

func main() {
	env := config.GetConfigVars()

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

	dg, err := discordgo.New("Bot " + env.DiscordAppSecret)
	if err != nil {
		log.Fatalf("Erro ao criar a sessão do Discord: %v", err)
	}

	err = dg.Open()
	if err != nil {
		log.Fatalf("Erro ao abrir a conexão: %v", err)
	}
	defer dg.Close()

	app := discordbot.NewApplication(discordbot.NewApplicationParams{
		RedisClient:     rdb,
		RedisChannel:    "send_to_discord",
		DiscordSession:  dg,
		DiscordChannel:  env.DiscordChannelID,
		MinioClient:     minioClient,
		MinioBucketName: env.MinioBucket,
	})

	app.Run()
}
