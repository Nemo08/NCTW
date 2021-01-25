package entity

import (
	"time"

	"github.com/google/uuid"
)

//Branch main model
type Branch struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"CreatedAt"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"UpdatedAt"`
	DeletedAt *time.Time `sql:"index" json:"-"`

	Description string //описание
	Name        string
}

//BranchResource is main model
type BranchResource struct {
	ID        uint64     `gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}
