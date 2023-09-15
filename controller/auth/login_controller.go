package auth

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()
var store = session.New()

func Login(c *fiber.Ctx) error {
	sess, _ := store.Get(c)
	message := sess.Get("message")
	sess.Delete("message")
	sess.Save()

	return c.Render("auth/login", fiber.Map{
		"message": message,
	})
}

func CheckLogin(c *fiber.Ctx) error {
	fmt.Println("hoho")
	email := c.FormValue("email")
	password := c.FormValue("password")

	var user model.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	if err == nil {
		fmt.Println(bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)))
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err == nil {
			sess, _ := store.Get(c)
			sess.Set("LoggedIn", user.ID)
			sess.Save()
			return c.Redirect("/dashboard")
		}
	}

	sess, _ := store.Get(c)
	sess.Set("message", "Invalid email or password")
	sess.Save()
	return c.Redirect("/login")
}
