package usecase

import (
	"github.com/google/uuid"

	ent "github.com/Nemo08/nctw/entity"
	int "github.com/Nemo08/nctw/interfaces"
)

type BranchUsecase interface {
	GetAllBranchs() ([]*ent.Branch, error)
	AddBranch(Branch ent.Branch) (ent.Branch, error)
	FindById(id uuid.UUID) (*ent.Branch, error)
	Find(q string) ([]*ent.Branch, error)
	UpdateBranch(Branch ent.Branch) (ent.Branch, error)
	DeleteBranchById(id uuid.UUID) error
}

type branchUsecase struct {
	repo int.BranchRepository
	log  LogInterface
}

func NewBranchUsecase(l LogInterface, r int.BranchRepository) *branchUsecase {
	return &branchUsecase{
		repo: r,
		log:  l,
	}
}

func (uc *branchUsecase) GetAllBranchs() ([]*ent.Branch, error) {
	uc.log.LogMessage("Get all branchs")

	branchs, err := uc.repo.GetAllBranchs()
	if err != nil {
		return nil, err
	}
	return branchs, nil
}

func (uc *branchUsecase) AddBranch(u ent.Branch) (ent.Branch, error) {
	uc.log.LogMessage("Add branch", u)
	return uc.repo.Store(u)
}

func (uc *branchUsecase) FindById(id uuid.UUID) (*ent.Branch, error) {
	uc.log.LogMessage("Find branch by id ", id)
	Branch, err := uc.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	return Branch, nil
}

func (uc *branchUsecase) UpdateBranch(u ent.Branch) (ent.Branch, error) {
	uc.log.LogMessage("Update branch", u)
	return uc.repo.UpdateBranch(u)
}

func (uc *branchUsecase) DeleteBranchById(id uuid.UUID) error {
	uc.log.LogMessage("Delete branch by id ", id)
	return uc.repo.DeleteBranchById(id)
}
