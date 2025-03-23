package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"

	_ "github.com/lib/pq"
)

type DBConfig struct {
	Username string `yaml:"dbUser"`
	Password string `yaml:"dbPassword"`
	Host     string `yaml:"dbHost"`
	Port     int    `yaml:"dbPort"`
	DbName   string `yaml:"dbName"`
	SSLMode  string `yaml:"sslmode"`
}

func NewDBConfig(cfg DBConfig) (*sqlx.DB, error) {
	db, err := connectDB(cfg)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectDB(cfg DBConfig) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DbName, cfg.SSLMode)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to the database")
	return db, nil
}

func CloseDB(db *sqlx.DB) error {
	if err := db.Close(); err != nil {
		return err
	}
	return nil
}
