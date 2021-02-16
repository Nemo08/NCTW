package validator

import (
	"github.com/go-playground/validator/v10"

	"github.com/Nemo08/NCTW/infrastructure/logger"
)

var validate *validator.Validate

func NewValidator() {
	logger.NewLogger().Info("Создаем валидатор")
	validate = validator.New()
}
