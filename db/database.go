package db

import (
	"database/sql"
	"fmt"
)

type Database struct {
	Conn *sql.DB
}

type Config struct {
	Host string
	Port int
	Username string
	Password string
	DbName string
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
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}
	return db, nil
}