package entity

import (
	"time"

	"github.com/google/uuid"
)

// Person main model
type Person struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"CreatedAt"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"UpdatedAt"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	Name       string `json:"name"`
	MiddleName string `json:"middle_name"`
	Surname    string `json:"surname"`

	Gender   uint       `json:"gender"` //1 - male, 2 - female
	Birthday *time.Time `json:"birthday"`
}
