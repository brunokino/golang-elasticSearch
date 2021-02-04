package main

import (
	"os"
	"strconv"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	db "gitlab.com/idoko/letterpress/db"
	"gitlab.com/idoko/letterpress/handler"
)

func main() {
	var dbPort int
	var err error
	var esClient *elasticsearch.Client
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	/*esConfig := elasticsearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://%s:%s", os.Getenv("ELASTIC_HOST"), os.Getenv("ELASTIC_PORT")),
		},
	}*/
	//esClient, err = elasticsearch.NewClient(esConfig)
	esClient, err = elasticsearch.NewDefaultClient()
	if err != nil {
		logger.Err(err).Msg("failed to create client")
	}
	port := os.Getenv("POSTGRES_PORT")
	if dbPort, err = strconv.Atoi(port); err != nil {
		logger.Err(err).Msg("failed to parse database port")
		os.Exit(1)
	}
	dbConfig := db.Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     dbPort,
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DbName:   os.Getenv("POSTGRES_DB"),

		ESClient: esClient,
		Logger:   logger,
	}
	logger.Info().Interface("config", &dbConfig).Msg("config:")
	dbInstance, err := db.Init(dbConfig)
	if err != nil {
		logger.Err(err).Msg("Connection failed")
		os.Exit(1)
	}
	logger.Info().Msg("Database connection established")

	h := handler.New(dbInstance, logger)
	router := gin.Default()
	rg := router.Group("/v1")
	h.Register(rg)
	router.Run(":8080")
}
