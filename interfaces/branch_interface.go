package interfaces

import (
	"github.com/google/uuid"

	ent "github.com/Nemo08/NCTW/entity"
)

type BranchRepository interface {
	Store(Branch ent.Branch) (*ent.Branch, error)
	GetAllBranchs() ([]*ent.Branch, error)
	FindByID(id uuid.UUID) (*ent.Branch, error)
	UpdateBranch(Branch ent.Branch) (*ent.Branch, error)
	DeleteBranchById(id uuid.UUID) error
}
