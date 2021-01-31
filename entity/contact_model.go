package entity

import (
	"github.com/google/uuid"
)

//Contact основная модель
type Contact struct {
	ID       uuid.UUID
	Position string //должность
}

//NewContact конструктор
func NewContact(position string) (Contact, error) {
	return Contact{
		ID:       uuid.New(),
		Position: position,
	}, nil
}
