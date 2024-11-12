package model

import (
	"database/sql"
	"time"
)

type AuditableEntity struct {
	CreatedAt time.Time     `db:"created_at"`
	CreatedBy string        `db:"created_by"`
	UpdatedAt sql.NullTime  `db:"updated_at"`
	UpdatedBy sql.NullInt64 `db:"updated_by"`
}
