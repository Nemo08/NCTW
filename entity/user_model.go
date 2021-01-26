package entity

import (
	"github.com/google/uuid"
)

// User базовая структура
type User struct {
	ID       uuid.UUID
	Login    string
	Password string
}
