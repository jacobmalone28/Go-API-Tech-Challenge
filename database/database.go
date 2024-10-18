package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/jacob-tech-challenge/config"
)

// OpenDBFunc is a function type that matches the signature of sql.Open
type OpenDBFunc func(driverName, dataSourceName string) (*sql.DB, error)

// Connect connects to the database using the provided OpenDBFunc
func Connect(cfg config.Config, openDB OpenDBFunc) (*sql.DB, error) {
	// create the data source name
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB_Host, cfg.DB_Port, cfg.DB_User, cfg.DB_Password, cfg.DB_Name)

	// connect to the database using the provided OpenDBFunc
	db, err := openDB("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	// check if the database is alive (by pinging it)
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}