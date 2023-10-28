package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
)

var ResultIDRole struct {
	ID uuid.UUID
}

type RegisterUser struct {
	Name     string `validate:"required,min=5,max=50"`
	Email    string `validate:"required,min=5,email,unique email"`
	Password string `validate:"required,min=3"`
	ImageUrl string
	Role     string
}

func Register(c *fiber.Ctx) error {
	return c.Render("auth/register", fiber.Map{
		"Errors": RegisterUser{},
	})
}

func CheckRegister(c *fiber.Ctx) error {
	database.DB.Table("roles").Select("id").Where("name = ?", "Pemohon").First(&ResultIDRole)
	registerUser := RegisterUser{
		Name:     c.FormValue("name"),
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
		ImageUrl: c.FormValue("image_url"),
	}

	registerValidator := &helper.Validator{
		Validator: validate,
	}

	if errs := registerValidator.Validate(registerUser); len(errs) > 0 {
		return showRegisterErrors(c, registerUser, errs)
	}

	hashedPassword, _ := helper.HashPassword(c.FormValue("password"))
	user := model.User{
		Name:     c.FormValue("name"),
		Email:    c.FormValue("email"),
		Password: hashedPassword,
		ImageUrl: c.FormValue("image_url"),
		RoleID:   ResultIDRole.ID,
		Business: &model.Business{},
	}

	result := database.DB.Create(&user)

	if result != nil {
		helper.SetSession(c, "messageSuccess", "Berhasil mendaftarkan akun baru! Silahkan login dengan akun baru tersebut!")
		return c.Redirect("/")
	}
	return nil
}

func showRegisterErrors(c *fiber.Ctx, oldInput RegisterUser, errs map[string]string) error {
	var errsStruct = RegisterUser{}
	if err := mapstructure.Decode(errs, &errsStruct); err != nil {
		panic(err)
	}
	return c.Render("auth/register", fiber.Map{
		"oldInput": oldInput,
		"Errors":   errsStruct,
	})
}
