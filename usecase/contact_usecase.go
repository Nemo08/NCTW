package usecase

import (
	"github.com/google/uuid"

	ent "github.com/Nemo08/NCTW/entity"
	log "github.com/Nemo08/NCTW/infrastructure/logger"
	int "github.com/Nemo08/NCTW/interfaces"
)

//ContactUsecase основная структура usecase
type ContactUsecase interface {
	GetAllContacts() ([]*ent.Contact, error)
	AddContact(Contact ent.Contact) (*ent.Contact, error)
	FindByID(id uuid.UUID) (*ent.Contact, error)
	Find(q string) ([]*ent.Contact, error)
	UpdateContact(Contact ent.Contact) (*ent.Contact, error)
	DeleteContactByID(id uuid.UUID) error
}

type contactUsecase struct {
	repo int.ContactRepository
	log  log.LogInterface
}

//NewContactUsecase создание объекта usecase для Contact
func NewContactUsecase(l log.LogInterface, r int.ContactRepository) *contactUsecase {
	return &contactUsecase{
		repo: r,
		log:  l,
	}
}

func (uc *contactUsecase) GetAllContacts() ([]*ent.Contact, error) {
	uc.log.LogMessage("Get all contacts")

	contacts, err := uc.repo.GetAllContacts()
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

func (uc *contactUsecase) AddContact(u ent.Contact) (*ent.Contact, error) {
	uc.log.LogMessage("Add contact", u)
	return uc.repo.Store(u)
}

func (uc *contactUsecase) FindByID(id uuid.UUID) (*ent.Contact, error) {
	uc.log.LogMessage("Find contact by id ", id)
	Contact, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return Contact, nil
}

func (uc *contactUsecase) Find(q string) ([]*ent.Contact, error) {
	uc.log.LogMessage("Find string info in contacts:", q)

	contacts, err := uc.repo.Find(q)
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

func (uc *contactUsecase) UpdateContact(u ent.Contact) (*ent.Contact, error) {
	uc.log.LogMessage("Update contact", u)
	return uc.repo.UpdateContact(u)
}

func (uc *contactUsecase) DeleteContactByID(id uuid.UUID) error {
	uc.log.LogMessage("Delete contact by id ", id)
	return uc.repo.DeleteContactByID(id)
}
