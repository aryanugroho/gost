package db

import (
	"time"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

type SqlxDBWrapper struct {
	*sqlx.DB
	timeout time.Duration
}

// NewSqlxDBWrapper wraps sql.DB with libs db wrapper ie. sqlx or gorm
func NewSqlxDBWrapper(db *sqlx.DB, driverName string, timeout time.Duration) *SqlxDBWrapper {
	return &SqlxDBWrapper{
		DB:      db,
		timeout: timeout,
	}
}

type GormDBWrapper struct {
	*gorm.DB
	timeout time.Duration
}

// NewDB wraps sql.DB with libs db wrapper ie. sqlx or gorm
func NewGormDBWrapper(db *gorm.DB, driverName string, timeout time.Duration) *GormDBWrapper {
	return &GormDBWrapper{
		DB:      db,
		timeout: timeout,
	}
}
