package usecase

import (
	"github.com/google/uuid"

	ent "github.com/Nemo08/NCTW/entity"
	repo "github.com/Nemo08/NCTW/infrastructure/repository"
)

type ContactUsecase interface {
	GetAllContacts() ([]*ent.Contact, error)
	AddContact(Contact ent.Contact) (ent.Contact, error)
	FindById(id uuid.UUID) (*ent.Contact, error)
	Find(q string) ([]*ent.Contact, error)
	UpdateContact(Contact ent.Contact) (ent.Contact, error)
	DeleteContactById(id uuid.UUID) error
}

type contactUsecase struct {
	repo int.ContactRepository
	log  LogInterface
}

func NewContactUsecase(l LogInterface, r int.ContactRepository) *contactUsecase {
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

func (uc *contactUsecase) AddContact(u ent.Contact) (ent.Contact, error) {
	uc.log.LogMessage("Add contact", u)
	return uc.repo.Store(u)
}

func (uc *contactUsecase) FindById(id uuid.UUID) (*ent.Contact, error) {
	uc.log.LogMessage("Find contact by id ", id)
	Contact, err := uc.repo.FindById(id)
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

func (uc *contactUsecase) UpdateContact(u ent.Contact) (ent.Contact, error) {
	uc.log.LogMessage("Update contact", u)
	return uc.repo.UpdateContact(u)
}

func (uc *contactUsecase) DeleteContactById(id uuid.UUID) error {
	uc.log.LogMessage("Delete contact by id ", id)
	return uc.repo.DeleteContactById(id)
}
