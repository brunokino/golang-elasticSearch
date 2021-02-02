package db

import (
	"database/sql"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type Database struct {
	Conn *sql.DB
	esClient *elasticsearch.Client
	Logger zerolog.Logger
}

type Config struct {
	Host string
	Port int
	Username string
	Password string
	DbName string
	ESClient *elasticsearch.Client
	Logger zerolog.Logger
}

func Init(cfg Config) (Database, error) {
	db := Database{}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DbName)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return db, err
	}

	db.Conn = conn
	db.esClient = cfg.ESClient
	db.Logger = cfg.Logger
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}
	return db, nil
}