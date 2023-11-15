package auth

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

var validate = validator.New()

func Login(c *fiber.Ctx) error {
	if c.Cookies("LoggedIn") != "" && c.Cookies("RoleUser") != "" {
		helper.SetSession(c, "LoggedIn", c.Cookies("LoggedIn"))
		helper.SetSession(c, "RoleUser", c.Cookies("RoleUser"))
		return c.Redirect("/pemohon/dashboard")
	}

	if role := helper.GetSession(c, "RoleUser"); role != "" {
		switch role {
		case "Pemohon":
			return c.Redirect("/pemohon/dashboard")
		case "Admin":
			return c.Redirect("/admin/dashboard")
		}
	}

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
	email := c.FormValue("email")
	password := c.FormValue("password")

	var user model.User
	err := database.DB.Preload("Role").Where("email = ?", email).First(&user).Error
	if err == nil {
		if user.IsActive == false {
			helper.SetSession(c, "message", "Akun belum aktif. Silahkan hubungi admin untuk aktivasi")
			return c.Redirect("/")
		}

		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err == nil {
			helper.SetSession(c, "LoggedIn", user.ID.String())
			helper.SetSession(c, "RoleUser", user.Role.Name)

			if c.FormValue("remember_me") == "on" {
				c.Cookie(&fiber.Cookie{
					Name:    "LoggedIn",
					Value:   user.ID.String(),
					Expires: time.Now().Add(24 * time.Hour),
				})

				c.Cookie(&fiber.Cookie{
					Name:    "RoleUser",
					Value:   user.Role.Name,
					Expires: time.Now().Add(24 * time.Hour),
				})
			}

			log := model.Log{
				UserID: user.ID,
			}
			database.DB.Create(&log)

			return c.Redirect(strings.ToLower(user.Role.Name) + "/dashboard")
		}
	}

	helper.SetSession(c, "message", "Invalid email or password")
	return c.Redirect("/")
}

func Logout(c *fiber.Ctx) error {
	helper.DestroySession(c)
	c.ClearCookie()
	return c.Redirect("/")
}
