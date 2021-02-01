package repository

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"gopkg.in/guregu/null.v4"

	ent "github.com/Nemo08/NCTW/entity"
	cfg "github.com/Nemo08/NCTW/infrastructure/config"
	log "github.com/Nemo08/NCTW/infrastructure/logger"
)

//DbContact стуктура для хранения Contact в базе
type DbContact struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `sql:"index"`

	Position null.String //должность
}

func db2contact(i DbContact) ent.Contact {
	return ent.Contact{
		ID:       i.ID,
		Position: i.Position,
	}
}

func contact2db(i ent.Contact) DbContact {
	return DbContact{
		ID:       i.ID,
		Position: i.Position,
	}
}

type contactRepositorySqlite struct {
	db  *gorm.DB
	log log.LogInterface
}

//NewContactRepositorySqlite создание объекта репозитория для Contact
func NewContactRepositorySqlite(l log.LogInterface, c cfg.ConfigInterface, db *gorm.DB) *contactRepositorySqlite {
	return &contactRepositorySqlite{
		db:  db,
		log: l,
	}
}

func (cts *contactRepositorySqlite) Store(contact ent.Contact) (*ent.Contact, error) {
	d := contact2db(contact)
	d.ID = uuid.New()

	errSlice := cts.db.Create(&d).GetErrors()
	if len(errSlice) != 0 {
		for _, err := range errSlice {
			cts.log.LogError("Error while contact create", err)
		}
		return &contact, errors.New("Error while contact create")
	}

	c := db2contact(d)
	return &c, nil
}

func (cts *contactRepositorySqlite) GetAllContacts() ([]*ent.Contact, error) {
	var c []*ent.Contact
	var d []*DbContact
	cts.db.Set("gorm:auto_preload", true).Find(&d)
	for _, r := range d {
		dbc := db2contact(*r)
		c = append(c, &dbc)
	}
	return c, nil
}

func (cts *contactRepositorySqlite) FindByID(id uuid.UUID) (*ent.Contact, error) {
	var d DbContact
	cts.db.Set("gorm:auto_preload", true).Where("id = ?", id).First(&d)
	c := db2contact(d)
	return &c, nil
}

func (cts *contactRepositorySqlite) Find(q string) ([]*ent.Contact, error) {
	var c []*ent.Contact
	var d []*DbContact
	cts.db.Set("gorm:auto_preload", true).Where("search_string LIKE ?", strings.ToLower("%"+q+"%")).Find(&d)
	for _, r := range d {
		dbc := db2contact(*r)
		c = append(c, &dbc)
	}
	return c, nil
}

func (cts *contactRepositorySqlite) UpdateContact(Contact ent.Contact) (*ent.Contact, error) {
	d := contact2db(Contact)
	//	c.SearchString = strings.ToLower(fmt.Sprintf("%v", c.Contact))
	cts.db.Set("gorm:auto_preload", true).Where("id = ?", d.ID).Save(&d)
	c := db2contact(d)
	return &c, nil
}

func (cts *contactRepositorySqlite) DeleteContactByID(id uuid.UUID) error {
	cts.db.Where("id = ?", id).Delete(DbContact{})
	return nil
}
