package helper

import (
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
	"gorm.io/gorm"
	"reflect"
)

type Validator struct {
	Validator *validator.Validate
}

type FileInput struct {
	Path        string
	Filename    string
	Ext         string
	ContentType string
	Size        int64
}

var validate = validator.New(validator.WithRequiredStructEnabled())
var english = en.New()
var uni = ut.New(english, english)
var trans, _ = uni.GetTranslator("en")

func (v Validator) Validate(data interface{}) map[string]string {
	_ = enTranslations.RegisterDefaultTranslations(validate, trans)

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("name")
	})

	validate.RegisterValidation("unique email", func(fl validator.FieldLevel) bool {
		var existingUser model.User
		result := database.DB.Where("email = ?", fl.Field().String()).First(&existingUser)
		if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
			return true
		}
		return false
	})

	validate.RegisterValidation("image_upload", func(fl validator.FieldLevel) bool {
		image := fl.Field().Interface().(FileInput)

		var acceptedImages = map[string]bool{
			"image/png":  true,
			"image/jpg":  true,
			"image/jpeg": true,
		}

		isValid := acceptedImages[image.ContentType]
		return isValid
	})

	imageUpload := "File {0} is not Valid"
	addTranslation("image_upload", imageUpload)

	val := make(map[string]string)

	err := validate.Struct(data)
	if err != nil {
		return translateError(err, val, trans)
	}

	return val
}

func translateError(err error, val map[string]string, trans ut.Translator) map[string]string {
	if err == nil {
		return nil
	}
	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(trans))
		val[e.StructField()] = translatedErr.Error()
	}
	return val
}

func addTranslation(tag string, errMessage string) {
	registerFn := func(ut ut.Translator) error {
		return ut.Add(tag, errMessage, false)
	}

	transFn := func(ut ut.Translator, fe validator.FieldError) string {
		param := fe.Param()
		tag := fe.Tag()

		t, err := ut.T(tag, fe.Field(), param)
		if err != nil {
			return fe.(error).Error()
		}
		return t
	}

	_ = validate.RegisterTranslation(tag, trans, registerFn, transFn)
}
