package user2

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/guregu/null.v4"
)

// User базовая структура
type User struct {
	ID           uuid.UUID
	Login        null.String
	PasswordHash null.String
	Email        null.String
}

//CreateHash создает хэш пароля
func CreateHash(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	return string(hash[:]), err
}

//NewUser конструктор
func NewUser(login, password, email null.String) (User, error) {
	stringHash, err := CreateHash(password.String)

	if err != nil {
		return User{}, err
	}

	return User{
		ID:           uuid.New(),
		Login:        login,
		PasswordHash: null.StringFrom(stringHash),
		Email:        email,
	}, err
}
