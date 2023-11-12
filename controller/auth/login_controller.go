package auth

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

var validate = validator.New()

func Login(c *fiber.Ctx) error {
	message := helper.GetSession(c, "message")
	helper.DeleteSession(c, "message")
	messageSuccess := helper.GetSession(c, "messageSuccess")
	helper.DeleteSession(c, "messageSuccess")

	return c.Render("auth/login", fiber.Map{
		"message":        message,
		"messageSuccess": messageSuccess,
	})
}

func CheckLogin(c *fiber.Ctx) error {
	fmt.Println("hoho")
	email := c.FormValue("email")
	password := c.FormValue("password")

	var user model.User
	err := database.DB.Preload("Role").Where("email = ?", email).First(&user).Error
	if err == nil {
		if user.IsActive == false {
			helper.SetSession(c, "message", "Akun belum aktif. Silahkan hubungi admin untuk aktivasi")
			return c.Redirect("/")
		}
		fmt.Println(bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)))
		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err == nil {
			helper.SetSession(c, "LoggedIn", user.ID.String())
			helper.SetSession(c, "RoleUser", user.Role.Name)

			return c.Redirect(strings.ToLower(user.Role.Name) + "/dashboard")
		}
	}

	helper.SetSession(c, "message", "Invalid email or password")
	return c.Redirect("/")
}

func Logout(c *fiber.Ctx) error {
	helper.DestroySession(c)
	return c.Redirect("/")
}
