package pemohon

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
	"path/filepath"
)

var validate = validator.New()

type UpdateProfileUser struct {
	Name  string       `validate:"required,min=5,max=50"`
	Image helper.Image `validate:"omitempty,image_upload"`
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
	img := helper.Image{}
	file, err := c.FormFile("profile_image")
	if err == nil {
		img = helper.Image{
			Path:        filepath.Dir(file.Filename),
			Filename:    filepath.Base(file.Filename),
			Ext:         filepath.Ext(file.Filename),
			ContentType: file.Header.Get("Content-Type"),
			Size:        file.Size,
		}
	}

	updateRegisterUser := UpdateProfileUser{
		Name:  c.FormValue("name"),
		Image: img,
	}

	registerValidator := &helper.Validator{
		Validator: validate,
	}

	if errs := registerValidator.Validate(updateRegisterUser); len(errs) > 0 {
		return showProfilePemohonErrors(c, updateRegisterUser, errs)
	}

	var user model.User
	err = database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	user.Name = c.FormValue("name")
	if newImage := UploadImageProfile(c); newImage != "" {
		user.ImageUrl = newImage
	}
	err = database.DB.Save(&user).Error
	helper.PanicIfError(err)

	helper.SetSession(c, "successMessage", "Profile Berhasil Diubah!")
	return c.Redirect("/pemohon/profile")
}

func UploadImageProfile(c *fiber.Ctx) string {
	filename := ""
	file, err := c.FormFile("profile_image")
	if err != nil {
		return filename
	}

	if file.Size != 0 {
		filename = "/uploads/profile/" + file.Filename
		err := c.SaveFile(file, "assets"+filename)
		helper.PanicIfError(err)
	}
	return filename
}

type MessageProfileUser struct {
	Name  string
	Image string
}

func showProfilePemohonErrors(c *fiber.Ctx, oldInput UpdateProfileUser, errs map[string]string) error {
	var user model.User
	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)
	fmt.Println(errs)
	var errsStruct = MessageProfileUser{}
	if err := mapstructure.Decode(errs, &errsStruct); err != nil {
		panic(err)
	}
	return c.Render("pemohon/profile/index", fiber.Map{
		"User":     user,
		"oldInput": oldInput,
		"Errors":   errsStruct,
	}, "layouts/pemohon")
}
