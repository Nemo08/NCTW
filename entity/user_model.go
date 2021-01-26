package entity

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User базовая структура
type User struct {
	ID           uuid.UUID
	Login        string
	PasswordHash string
	Email        string
}

func NewUser(login, password, email string) (User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	return User{
		ID:           uuid.New(),
		Login:        login,
		PasswordHash: string(hash[:]),
		Email:        email,
	}, err
}
