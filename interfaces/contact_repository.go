package interfaces

import (
	"github.com/google/uuid"

	ent "github.com/Nemo08/nctw/entity"
)

type ContactRepository interface {
	Store(Contact ent.Contact) (ent.Contact, error)
	GetAllContacts() ([]*ent.Contact, error)
	FindById(id uuid.UUID) (*ent.Contact, error)
	Find(q string) ([]*ent.Contact, error)
	UpdateContact(Contact ent.Contact) (ent.Contact, error)
	DeleteContactById(id uuid.UUID) error
}
