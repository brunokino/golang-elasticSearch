package main

import (
	"github.com/bxcodec/faker/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.com/idoko/letterpress/db"
	"gitlab.com/idoko/letterpress/models"
	"os"
)

// maybe run this as a binary service in docker and make api depend on it?
func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	dbConfig := db.Config{
		// hard-coding connection info since we only have to seed on the dev machine
		Host: "localhost",
		Port: 5432,
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DbName: os.Getenv("POSTGRES_DB"),
	}

	dbInstance, err := db.Init(dbConfig)
	if err != nil {
		logger.Err(err).Msg("Connection failed")
		os.Exit(1)
	}

	for i := 0; i < 20; i++ {
		post := &models.Post{
			Title: faker.Sentence(),
			Body: faker.Paragraph(),
		}
		err = dbInstance.SavePost(post)
		if err != nil {
			log.Err(err).Msg("failed to save record")
		}
	}
}