package usecase

import (
	"github.com/google/uuid"

	ent "github.com/Nemo08/NCTW/entity"
	log "github.com/Nemo08/NCTW/infrastructure/logger"
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

type ContactUsecaseStruct struct {
	repo ContactRepository
	log  log.LogInterface
}

//NewContactUsecase создание объекта usecase для Contact
func NewContactUsecase(l log.LogInterface, r ContactRepository) *ContactUsecaseStruct {
	return &ContactUsecaseStruct{
		repo: r,
		log:  l,
	}
}

func (uc *ContactUsecaseStruct) GetAllContacts() ([]*ent.Contact, error) {
	uc.log.LogMessage("Get all contacts")

	contacts, err := uc.repo.GetAllContacts()
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

func (uc *ContactUsecaseStruct) AddContact(u ent.Contact) (*ent.Contact, error) {
	uc.log.LogMessage("Add contact", u)
	return uc.repo.Store(u)
}

func (uc *ContactUsecaseStruct) FindByID(id uuid.UUID) (*ent.Contact, error) {
	uc.log.LogMessage("Find contact by id ", id)
	Contact, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return Contact, nil
}

func (uc *ContactUsecaseStruct) Find(q string) ([]*ent.Contact, error) {
	uc.log.LogMessage("Find string info in contacts:", q)

	contacts, err := uc.repo.Find(q)
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

func (uc *ContactUsecaseStruct) UpdateContact(u ent.Contact) (*ent.Contact, error) {
	uc.log.LogMessage("Update contact", u)
	return uc.repo.UpdateContact(u)
}

func (uc *ContactUsecaseStruct) DeleteContactByID(id uuid.UUID) error {
	uc.log.LogMessage("Delete contact by id ", id)
	return uc.repo.DeleteContactByID(id)
}
