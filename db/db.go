package db

import (
	"fmt"
	"log"
	"wallet-manager/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDB(cfg *config.DatabaseConfig) *sqlx.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}
