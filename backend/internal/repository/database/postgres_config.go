package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

type DBSettings struct {
	DBUser     string `yaml:"dbUser"`
	DBPassword string `yaml:"dbPassword"`
	DBHost     string `yaml:"dbHost"`
	DBPort     int    `yaml:"dbPort"`
	DBName     string `yaml:"dbName"`
	SSLMode    string `yaml:"sslmode"`
}

func NewDBConfig(cfg DBSettings) (*sqlx.DB, error) {
	db, err := connectDB(cfg)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectDB(cfg DBSettings) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.SSLMode)
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
