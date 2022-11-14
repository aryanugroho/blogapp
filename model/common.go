package model

import (
	"database/sql"
	"time"
)

type AuditableEntity struct {
	CreatedAt *time.Time     `gorm:"created_at"`
	CreatedBy *int64         `gorm:"created_by"`
	UpdatedAt *sql.NullTime  `gorm:"updated_at"`
	UpdatedBy *sql.NullInt64 `gorm:"updated_by"`
}
