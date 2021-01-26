package interfaces

import (
	"github.com/google/uuid"

	ent "github.com/Nemo08/NCTW/entity"
)

type ContactRepository interface {
	Store(Contact ent.Contact) (ent.Contact, error)
	GetAllContacts() ([]*ent.Contact, error)
	FindByID(id uuid.UUID) (*ent.Contact, error)
	Find(q string) ([]*ent.Contact, error)
	UpdateContact(Contact ent.Contact) (ent.Contact, error)
	DeleteContactByID(id uuid.UUID) error
}
