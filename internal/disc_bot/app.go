package discbot

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
	redismodel "github.com/v-venes/lol-feed-journal/pkg/models/redis"
)

type Application struct {
	RedisClient     *redis.Client
	RedisChannel    string
	DiscordSession  *discordgo.Session
	DiscordChannel  string
	MinioClient     *minio.Client
	MinioBucketName string
}

type NewApplicationParams struct {
	RedisChannel    string
	DiscordChannel  string
	RedisClient     *redis.Client
	DiscordSession  *discordgo.Session
	MinioClient     *minio.Client
	MinioBucketName string
}

func NewApplication(params NewApplicationParams) *Application {

	return &Application{
		RedisClient:    params.RedisClient,
		RedisChannel:   params.RedisChannel,
		DiscordSession: params.DiscordSession,
		DiscordChannel: params.DiscordChannel,
		MinioClient:    params.MinioClient,
	}
}

func (a *Application) Run() {
	ctx := context.Background()

	pubsub := a.RedisClient.Subscribe(ctx, a.RedisChannel)
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Fatal(err)
		}

		var message redismodel.SendToDiscordPayload

		err = json.Unmarshal([]byte(msg.Payload), &message)

		if err != nil {
			log.Fatal(err)
		}

		err = a.getJournalAndSend(message.JournalPath)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func (a *Application) getJournalAndSend(journalPath string) error {
	downloadPath := "/tmp/journal.png"

	err := a.MinioClient.FGetObject(context.Background(), a.MinioBucketName, journalPath, downloadPath, minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	file, err := os.Open(downloadPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = a.DiscordSession.ChannelMessageSendComplex(a.DiscordChannel, &discordgo.MessageSend{
		Content: "@everyone Jornal do Feed do dia anterior",
		Files: []*discordgo.File{
			{
				Name:   downloadPath,
				Reader: file,
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}
