package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/v-venes/lol-feed-journal/internal/config"
	domain "github.com/v-venes/lol-feed-journal/internal/domain/repository"
	"github.com/v-venes/lol-feed-journal/internal/repository"
)

func init() {
	if os.Getenv("GO_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("[ERROR] %s", err.Error())
		}
	}
}

func main() {
	env := config.GetConfigVars()
	db, err := repository.NewPostgresDB(repository.DatabaseConfig{
		DSN: fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			env.PostgresHost,
			env.PostgresUser,
			env.PostgresPassword,
			env.PostgresDB,
			env.PostgresPort,
		),
	})
	if err != nil {
		log.Fatalf("[ERROR] %s", err.Error())
	}

	err = db.Migrator().AutoMigrate(&domain.Player{}, &domain.Match{}, &domain.MatchPlayer{})
	if err != nil {
		log.Fatalf("[ERROR] Migration Error: %s", err.Error())
	}

}
