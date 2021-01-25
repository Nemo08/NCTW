package entity

import (
	"time"

	"github.com/google/uuid"
)

//Contact main model
type Contact struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"CreatedAt"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"UpdatedAt"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	Position string //должность
}
