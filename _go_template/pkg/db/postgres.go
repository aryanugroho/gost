package db

import (
	"database/sql"
	"fmt"
	"time"
)

func PostgresDB(config Config, timeout time.Duration) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.Name)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(config.MaxOpen)
	db.SetMaxIdleConns(config.MaxIdle)
	db.SetConnMaxLifetime(time.Duration(config.MaxLifetime) * time.Minute)

	// Ping the database to check the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
