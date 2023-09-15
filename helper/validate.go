package helper

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	Validator *validator.Validate
}

var validate = validator.New()

func (v Validator) Validate(data interface{}) map[string]string {
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
