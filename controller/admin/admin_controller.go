package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
)

func ListAdmin(c *fiber.Ctx) error {
	message := helper.GetSession(c, "successMessage")
	helper.DeleteSession(c, "successMessage")
	var users []model.User
	var user model.User

	var role model.Role
	err := database.DB.First(&role, "name = ?", "Admin").Error
	helper.PanicIfError(err)

	err = database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	err = database.DB.Where("role_id = ?", role.ID).Where("id != ?", helper.GetSession(c, "LoggedIn")).Find(&users).Error
	helper.PanicIfError(err)

	return c.Render("admin/user/list-admin", fiber.Map{
		"User":    user,
		"Users":   users,
		"message": message,
	}, "layouts/admin")
}

func AddAdmin(c *fiber.Ctx) error {
	var user model.User

	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)
	return c.Render("admin/user/add-admin", fiber.Map{
		"User": user,
	}, "layouts/admin")
}

type AddAdminUser struct {
	Name     string           `validate:"required,min=5,max=50"`
	Email    string           `validate:"required,min=5,email,unique email"`
	Password string           `validate:"required,min=3"`
	ImageUrl helper.FileInput `validate:"omitempty,image_upload"`
	Role     string
}

func StoreAdmin(c *fiber.Ctx) error {
	imgProfile, updateImgProfile := helper.CheckInputFile(c, "image_url")

	var role model.Role
	err := database.DB.First(&role, "name = ?", "Admin").Error
	helper.PanicIfError(err)

	registerUser := AddAdminUser{
		Name:     c.FormValue("name"),
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
		ImageUrl: imgProfile,
	}

	registerValidator := &helper.Validator{
		Validator: validate,
	}

	if errs := registerValidator.Validate(registerUser); len(errs) > 0 {
		return showAddAdminErrors(c, registerUser, errs)
	}

	path := ""
	if updateImgProfile {
		if pathLogo, ok := helper.UploadFile(c, "image_url", "profile/user"); ok {
			path = pathLogo
		}
	}

	hashedPassword, _ := helper.HashPassword(c.FormValue("password"))
	user := model.User{
		Name:     c.FormValue("name"),
		Email:    c.FormValue("email"),
		Password: hashedPassword,
		ImageUrl: path,
		RoleID:   role.ID,
		Business: &model.Business{},
	}

	result := database.DB.Create(&user)

	if result != nil {
		helper.SetSession(c, "successMessage", "Berhasil menambahkan akun Admin!")
		return c.Redirect("/admin/user/list-admin")
	}
	return nil
}

func EditAdmin(c *fiber.Ctx) error {
	id := c.Params("userId")
	var user model.User
	var userEdit model.User

	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	err = database.DB.First(&userEdit, "id = ?", id).Error
	helper.PanicIfError(err)

	return c.Render("admin/user/edit-admin", fiber.Map{
		"User":     user,
		"UserEdit": userEdit,
	}, "layouts/admin")
}

type EditAdminUser struct {
	Name     string           `validate:"required,min=5,max=50"`
	Email    string           `validate:"required,min=5,email"`
	Password string           `validate:"omitempty,min=3"`
	ImageUrl helper.FileInput `validate:"omitempty,image_upload"`
	Role     string
}

func UpdateUserAdmin(c *fiber.Ctx) error {
	id := c.Params("userId")
	var user model.User

	err := database.DB.First(&user, "id = ?", id).Error
	helper.PanicIfError(err)

	imgProfile, updateImgProfile := helper.CheckInputFile(c, "image_url")

	editUser := EditAdminUser{
		Name:     c.FormValue("name"),
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
		ImageUrl: imgProfile,
	}

	registerValidator := &helper.Validator{
		Validator: validate,
	}

	if errs := registerValidator.Validate(editUser); len(errs) > 0 {
		return showEditAdminErrors(c, editUser, errs)
	}

	user.Name = c.FormValue("name")
	user.Email = c.FormValue("email")
	if pass := c.FormValue("password"); pass != "" {
		hashedPassword, _ := helper.HashPassword(pass)
		user.Password = hashedPassword
	}

	if updateImgProfile {
		if path, ok := helper.UploadFile(c, "image_url", "profile/user"); ok {
			user.ImageUrl = path
		}
	}

	err = database.DB.Save(&user).Error
	helper.PanicIfError(err)

	helper.SetSession(c, "successMessage", "Berhasil mengedit akun Admin!")
	return c.Redirect("/admin/user/list-admin")
}

func DeleteAdmin(c *fiber.Ctx) error {
	id := c.Params("userId")
	var user model.User

	err := database.DB.First(&user, "id = ?", id).Error
	helper.PanicIfError(err)

	//_ = os.Remove(path.Join("assets", user.BrandLogo))
	err = database.DB.Delete(&user).Error
	helper.PanicIfError(err)

	helper.SetSession(c, "successMessage", "Berhasil menghapus akun Admin!")
	return c.Redirect("/admin/user/list-admin")
}

type MessageAddEditAdmin struct {
	Name     string
	Email    string
	Password string
	ImageUrl string
	Role     string
}

func showAddAdminErrors(c *fiber.Ctx, oldInput AddAdminUser, errs map[string]string) error {
	var errsStruct = MessageAddEditAdmin{}
	if err := mapstructure.Decode(errs, &errsStruct); err != nil {
		panic(err)
	}
	return c.Render("admin/user/add-admin", fiber.Map{
		"oldInput": oldInput,
		"Errors":   errsStruct,
	}, "layouts/admin")
}

func showEditAdminErrors(c *fiber.Ctx, oldInput EditAdminUser, errs map[string]string) error {
	var errsStruct = MessageAddEditAdmin{}
	if err := mapstructure.Decode(errs, &errsStruct); err != nil {
		panic(err)
	}
	return c.Render("admin/user/edit-admin", fiber.Map{
		"oldInput": oldInput,
		"Errors":   errsStruct,
	}, "layouts/admin")
}
