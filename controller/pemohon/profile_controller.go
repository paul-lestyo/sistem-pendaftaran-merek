package pemohon

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
)

var validate = validator.New()

type UpdateProfileUser struct {
	Name string `validate:"required,min=5,max=50"`
}

func ProfilePemohon(c *fiber.Ctx) error {
	message := helper.GetSession(c, "successMessage")
	helper.DeleteSession(c, "successMessage")
	var user model.User

	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)
	return c.Render("pemohon/profile/index", fiber.Map{
		"User":    user,
		"message": message,
	}, "layouts/pemohon")
}

func UpdatePemohon(c *fiber.Ctx) error {
	updateRegisterUser := UpdateProfileUser{
		Name: c.FormValue("name"),
	}

	registerValidator := &helper.Validator{
		Validator: validate,
	}

	if errs := registerValidator.Validate(updateRegisterUser); len(errs) > 0 {
		return showProfilePemohonErrors(c, updateRegisterUser, errs)
	}

	var user model.User
	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	user.Name = c.FormValue("name")
	err = database.DB.Save(&user).Error
	helper.PanicIfError(err)

	helper.SetSession(c, "successMessage", "Profile Berhasil Diubah!")
	return c.Redirect("/pemohon/profile")
}

func showProfilePemohonErrors(c *fiber.Ctx, oldInput UpdateProfileUser, errs map[string]string) error {
	var errsStruct = UpdateProfileUser{}
	if err := mapstructure.Decode(errs, &errsStruct); err != nil {
		panic(err)
	}
	return c.Render("pemohon/profile/index", fiber.Map{
		"oldInput": oldInput,
		"Errors":   errsStruct,
	}, "layouts/pemohon")
}
