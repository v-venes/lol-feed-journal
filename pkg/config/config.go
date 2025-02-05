package config

import (
	"log"

	"github.com/Netflix/go-env"
)

type Config struct {
	RiotBasePath     string `env:"RIOT_API_BASE_PATH"`
	RiotKey          string `env:"RIOT_API_KEY"`
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresDB       string `env:"POSTGRES_DB"`
	RedisHost        string `env:"REDIS_HOST"`
	RedisChannel     string `env:"REDIS_CHANNEL"`
	MinioHost        string `env:"MINIO_HOST"`
	MinioKey         string `env:"MINIO_KEY"`
	MinioSecret      string `env:"MINIO_SECRET"`
	DiscordAppID     string `env:"DISCORD_APP_ID"`
	DiscordAppSecret string `env:"DISCORD_APP_KEY"`
	DiscordChannelID string `env:"DISCORD_CHANNEL_ID"`
}

func GetConfigVars() *Config {
	var config Config
	_, err := env.UnmarshalFromEnviron(&config)

	if err != nil {
		log.Fatal(err)
	}

	return &config
}
