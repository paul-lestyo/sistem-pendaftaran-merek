package admin

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
)

var validate = validator.New()

type UpdateProfileUser struct {
	Name  string           `validate:"required,min=5,max=50"`
	Image helper.FileInput `validate:"omitempty,image_upload"`
}

func ProfileAdmin(c *fiber.Ctx) error {
	message := helper.GetSession(c, "successMessage")
	helper.DeleteSession(c, "successMessage")
	var user model.User

	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)
	return c.Render("admin/profile", fiber.Map{
		"User":    user,
		"message": message,
	}, "layouts/admin")
}

func UpdateProfileAdmin(c *fiber.Ctx) error {
	img, updateImg := helper.CheckInputFile(c, "profile_image")
	updateRegisterUser := UpdateProfileUser{
		Name:  c.FormValue("name"),
		Image: img,
	}

	registerValidator := &helper.Validator{
		Validator: validate,
	}

	if errs := registerValidator.Validate(updateRegisterUser); len(errs) > 0 {
		return showProfileAdminErrors(c, updateRegisterUser, errs)
	}

	var user model.User
	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	user.Name = c.FormValue("name")
	if updateImg {
		if path, ok := helper.UploadFile(c, "profile_image", "profile/user"); ok {
			user.ImageUrl = path
		}
	}
	err = database.DB.Save(&user).Error
	helper.PanicIfError(err)

	helper.SetSession(c, "successMessage", "Profile Berhasil Diubah!")
	return c.Redirect("/admin/profile")
}

type MessageProfileUser struct {
	Name  string
	Image string
}

func showProfileAdminErrors(c *fiber.Ctx, oldInput UpdateProfileUser, errs map[string]string) error {
	var user model.User
	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)
	fmt.Println(errs)
	var errsStruct = MessageProfileUser{}
	if err := mapstructure.Decode(errs, &errsStruct); err != nil {
		panic(err)
	}
	return c.Render("admin/profile", fiber.Map{
		"User":     user,
		"oldInput": oldInput,
		"Errors":   errsStruct,
	}, "layouts/admin")
}
