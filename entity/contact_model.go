package entity

import (
	"github.com/google/uuid"
	"gopkg.in/guregu/null.v4"
)

//Contact основная модель
type Contact struct {
	ID       uuid.UUID
	Position null.String //должность
}

//NewContact конструктор
func NewContact(position null.String) (Contact, error) {
	return Contact{
		ID:       uuid.New(),
		Position: position,
	}, nil
}
