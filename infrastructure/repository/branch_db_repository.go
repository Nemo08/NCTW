package repository

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	ent "github.com/Nemo08/NCTW/entity"
	cfg "github.com/Nemo08/NCTW/infrastructure/config"
	use "github.com/Nemo08/NCTW/usecase"
)

type DbBranch struct {
	ent.Branch
}

type BranchRepositorySqlite struct {
	db  *gorm.DB
	log use.LogInterface
}

func NewBranchRepositorySqlite(l use.LogInterface, c cfg.ConfigInterface, db *gorm.DB) *BranchRepositorySqlite {
	return &BranchRepositorySqlite{
		db:  db,
		log: l,
	}
}

func (cts *BranchRepositorySqlite) Store(Branch ent.Branch) (ent.Branch, error) {
	var c DbBranch
	//	var p ent.Person
	c.Branch = Branch
	a := uuid.New()

	c.ID = a

	err_slice := cts.db.Create(&c).GetErrors()
	if len(err_slice) != 0 {
		for _, err := range err_slice {
			cts.log.LogError("Error while Branch create", err)
		}
		return c.Branch, errors.New("Error while Branch create")
	}
	return c.Branch, nil
}

func (cts *BranchRepositorySqlite) GetAllBranchs() ([]*ent.Branch, error) {
	var Branchs []*ent.Branch
	var dbBranchs []*DbBranch
	cts.db.Set("gorm:auto_preload", true).Find(&dbBranchs)
	for _, c := range dbBranchs {
		Branchs = append(Branchs, &c.Branch)
	}
	return Branchs, nil
}

func (cts *BranchRepositorySqlite) FindById(id uuid.UUID) (*ent.Branch, error) {
	var c DbBranch
	cts.db.Set("gorm:auto_preload", true).Where("id = ?", id).First(&c)
	return &c.Branch, nil
}

func (cts *BranchRepositorySqlite) Find(q string) ([]*ent.Branch, error) {
	var Branchs []*ent.Branch
	var dbBranchs []*DbBranch
	cts.db.Set("gorm:auto_preload", true).Where("search_string LIKE ?", strings.ToLower("%"+q+"%")).Find(&dbBranchs)
	for _, c := range dbBranchs {
		Branchs = append(Branchs, &c.Branch)
	}
	return Branchs, nil
}

func (cts *BranchRepositorySqlite) UpdateBranch(Branch ent.Branch) (ent.Branch, error) {
	var c DbBranch
	c.Branch = Branch
	cts.db.Set("gorm:auto_preload", true).Where("id = ?", c.Branch.ID).Save(&c)
	return c.Branch, nil
}

func (cts *BranchRepositorySqlite) DeleteBranchById(id uuid.UUID) error {
	cts.db.Where("id = ?", id).Delete(DbBranch{})
	return nil
}
