package imagegenerator

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang/freetype/truetype"
	feedjournal "github.com/v-venes/feed_journal/pkg/models/feed_journal"
	redismodel "github.com/v-venes/feed_journal/pkg/models/redis"
	"github.com/v-venes/feed_journal/pkg/repositories"
	"github.com/v-venes/feed_journal/pkg/services"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/minio/minio-go/v7"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

type Application struct {
	PlayerRepository *repositories.PlayerRepository
	MatchRepository  *repositories.MatchRepository
	LeagueService    *services.LeagueService
	RedisClient      *redis.Client
	RedisChannel     string
	MinioClient      *minio.Client
}

type NewApplicationParams struct {
	PlayerRepository *repositories.PlayerRepository
	MatchRepository  *repositories.MatchRepository
	LeagueService    *services.LeagueService
	RedisClient      *redis.Client
	RedisChannel     string
	MinioClient      *minio.Client
}

type TemplateDimensions struct {
	ChampionImageYOffset        int
	ChampionImageDistanceOffset int
	ChampionLeftImageX          int
	ChampionRightImageX         int
	ChampionImageSize           int

	DateYOffset float64
	DateXOffset float64

	PlayerNameYOffset        float64
	PlayerKDAYOffset         float64
	PlayerLeftInfoXOffset    float64
	PlayerRightInfoXOffset   float64
	PlayerInfoDistanceOffset float64
}

func NewApplication(params NewApplicationParams) *Application {
	return &Application{
		PlayerRepository: params.PlayerRepository,
		MatchRepository:  params.MatchRepository,
		LeagueService:    params.LeagueService,
		RedisClient:      params.RedisClient,
		RedisChannel:     params.RedisChannel,
		MinioClient:      params.MinioClient,
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

		var message redismodel.GenerateImagePayload

		err = json.Unmarshal([]byte(msg.Payload), &message)

		if err != nil {
			log.Fatal(err)
		}

		journalPath, err := a.startImageGeneration(message.MatchDate)

		if err != nil {
			log.Fatal(err)
		}

		discordMessage := redismodel.SendToDiscordPayload{JournalPath: journalPath}
		jsonMsg, _ := json.Marshal(discordMessage)
		err = a.RedisClient.Publish(ctx, "send_to_discord", jsonMsg).Err()

		if err != nil {
			log.Fatal(err)
		}
	}
}

func (a *Application) startImageGeneration(matchDate time.Time) (string, error) {

	journalInfo, err := a.getJournalInfo(matchDate)

	if err != nil {
		return "", err
	}

	downloadPath := "/tmp/template.png"
	err = a.MinioClient.FGetObject(context.Background(), "feedjournal", "template/feed_journal_template.png", downloadPath, minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}

	// TODO: refatorar essa parte, fazendo tudo junto só pra testar

	dimensions := TemplateDimensions{
		ChampionImageYOffset:        162,
		ChampionImageDistanceOffset: 53,
		ChampionLeftImageX:          206,
		ChampionRightImageX:         538,
		ChampionImageSize:           38,

		DateXOffset: 418,
		DateYOffset: 134,

		PlayerNameYOffset:        180,
		PlayerKDAYOffset:         200,
		PlayerLeftInfoXOffset:    275.4,
		PlayerRightInfoXOffset:   608.4,
		PlayerInfoDistanceOffset: 51.7,
	}
	im1, err := gg.LoadPNG(downloadPath)
	if err != nil {
		return "", err
	}

	s1 := im1.Bounds().Size()
	dc := gg.NewContext(s1.X, s1.Y)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.DrawImage(im1, 0, 0)
	dc.SetRGB(0, 0, 0)
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return "", err
	}
	face := truetype.NewFace(font, &truetype.Options{
		Size: 16,
	})
	dc.SetFontFace(face)

	dc.DrawString(matchDate.Format("02/01/2006"), dimensions.DateXOffset, dimensions.DateYOffset)

	for index, player := range journalInfo.StoredPlayers {
		dc.DrawString(player.SummonerName, dimensions.PlayerLeftInfoXOffset, dimensions.PlayerNameYOffset+(dimensions.PlayerInfoDistanceOffset*float64(index)))
		kdaText := fmt.Sprintf("%d/%d/%d", player.Kills, player.Deaths, player.Assists)
		dc.DrawString(kdaText, dimensions.PlayerLeftInfoXOffset, dimensions.PlayerKDAYOffset+(dimensions.PlayerInfoDistanceOffset*float64(index)))

		imagePath, err := a.LeagueService.DownloadChampionSquareImage(player.ChampionName)
		if err != nil {
			return "", err
		}

		img, err := gg.LoadPNG(*imagePath)
		championImg := imaging.Resize(img, dimensions.ChampionImageSize, dimensions.ChampionImageSize, imaging.Lanczos)
		if err != nil {
			return "", err
		}

		dc.DrawImage(championImg, dimensions.ChampionLeftImageX, dimensions.ChampionImageYOffset+(dimensions.ChampionImageDistanceOffset*index))
	}

	for index, player := range journalInfo.RandomPlayers {
		dc.DrawString(player.SummonerName, dimensions.PlayerRightInfoXOffset, dimensions.PlayerNameYOffset+(dimensions.PlayerInfoDistanceOffset*float64(index)))
		kdaText := fmt.Sprintf("%d/%d/%d", player.Kills, player.Deaths, player.Assists)
		dc.DrawString(kdaText, float64(dimensions.PlayerRightInfoXOffset), dimensions.PlayerKDAYOffset+(dimensions.PlayerInfoDistanceOffset*float64(index)))

		imagePath, err := a.LeagueService.DownloadChampionSquareImage(player.ChampionName)
		if err != nil {
			return "", err
		}

		img, err := gg.LoadPNG(*imagePath)
		championImg := imaging.Resize(img, dimensions.ChampionImageSize, dimensions.ChampionImageSize, imaging.Lanczos)
		if err != nil {
			return "", err
		}

		dc.DrawImage(championImg, dimensions.ChampionRightImageX, dimensions.ChampionImageYOffset+(dimensions.ChampionImageDistanceOffset*index))
	}

	dc.SavePNG("./out.png")

	file, err := os.Open("./out.png")
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("journals/journal-%s.png", matchDate.Format("2006-01-02"))
	_, err = a.MinioClient.PutObject(context.Background(), "feedjournal", fileName, file, fileInfo.Size(), minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (a *Application) getJournalInfo(matchDate time.Time) (*feedjournal.Journal, error) {
	storedPlayers, err := a.MatchRepository.GetTop10StoredPlayersMatches(matchDate)

	if err != nil {
		return nil, err
	}

	randomPlayers, err := a.MatchRepository.GetTop10RandomPlayersMatches(matchDate)

	if err != nil {
		return nil, err
	}

	journal := &feedjournal.Journal{
		StoredPlayers: storedPlayers,
		RandomPlayers: randomPlayers,
		MatchDate:     matchDate,
	}

	return journal, nil
}
