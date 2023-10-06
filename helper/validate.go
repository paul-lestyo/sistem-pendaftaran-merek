package helper

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
	"gorm.io/gorm"
)

type Validator struct {
	Validator *validator.Validate
}

var validate = validator.New()

func (v Validator) Validate(data interface{}) map[string]string {

	validate.RegisterValidation("unique email", func(fl validator.FieldLevel) bool {
		var existingUser model.User
		result := database.DB.Where("email = ?", fl.Field().String()).First(&existingUser)
		if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
			return true
		}
		return false
	})

	validationErrors := make(map[string]string)

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			validationErrors[err.Field()] = fmt.Sprintf("%s field is %s %s", err.Field(), err.Tag(), err.Param())
		}
	}

	return validationErrors
}
