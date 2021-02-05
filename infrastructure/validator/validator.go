package validator

import (
	"github.com/go-playground/validator/v10"

	log "github.com/Nemo08/NCTW/infrastructure/logger"
)

var validate *validator.Validate

func NewValidator() {
	log.LogMessage("Создаем валидатор")
	validate = validator.New()
}
