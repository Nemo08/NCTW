package repository

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	ent "github.com/Nemo08/NCTW/entity"
	cfg "github.com/Nemo08/NCTW/infrastructure/config"
	log "github.com/Nemo08/NCTW/infrastructure/logger"
)

type DbContact struct {
	ent.Contact
}

type contactRepositorySqlite struct {
	db  *gorm.DB
	log log.LogInterface
}

func NewContactRepositorySqlite(l log.LogInterface, c cfg.ConfigInterface, db *gorm.DB) *contactRepositorySqlite {
	return &contactRepositorySqlite{
		db:  db,
		log: l,
	}
}

func (cts *contactRepositorySqlite) Store(Contact ent.Contact) (ent.Contact, error) {
	var c DbContact
	c.Contact = Contact
	a := uuid.New()

	c.ID = a
	err_slice := cts.db.Create(&c).GetErrors()
	if len(err_slice) != 0 {
		for _, err := range err_slice {
			cts.log.LogError("Error while contact create", err)
		}
		return c.Contact, errors.New("Error while contact create")
	}
	return c.Contact, nil
}

func (cts *contactRepositorySqlite) GetAllContacts() ([]*ent.Contact, error) {
	var contacts []*ent.Contact
	var dbcontacts []*DbContact
	cts.db.Set("gorm:auto_preload", true).Find(&dbcontacts)
	for _, c := range dbcontacts {
		contacts = append(contacts, &c.Contact)
	}
	return contacts, nil
}

func (cts *contactRepositorySqlite) FindById(id uuid.UUID) (*ent.Contact, error) {
	var c DbContact
	cts.db.Set("gorm:auto_preload", true).Where("id = ?", id).First(&c)
	return &c.Contact, nil
}

func (cts *contactRepositorySqlite) Find(q string) ([]*ent.Contact, error) {
	var contacts []*ent.Contact
	var dbcontacts []*DbContact
	cts.db.Set("gorm:auto_preload", true).Where("search_string LIKE ?", strings.ToLower("%"+q+"%")).Find(&dbcontacts)
	for _, c := range dbcontacts {
		contacts = append(contacts, &c.Contact)
	}
	return contacts, nil
}

func (cts *contactRepositorySqlite) UpdateContact(Contact ent.Contact) (ent.Contact, error) {
	var c DbContact
	c.Contact = Contact
	//	c.SearchString = strings.ToLower(fmt.Sprintf("%v", c.Contact))
	cts.db.Set("gorm:auto_preload", true).Where("id = ?", c.Contact.ID).Save(&c)
	return c.Contact, nil
}

func (cts *contactRepositorySqlite) DeleteContactById(id uuid.UUID) error {
	cts.db.Where("id = ?", id).Delete(DbContact{})
	return nil
}
