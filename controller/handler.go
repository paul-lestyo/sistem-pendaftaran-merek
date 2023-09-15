package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
)

func GetAllUsers(c *fiber.Ctx) error {
	db := database.DB
	var users []model.User

	db.Find(&users)
	return c.Render("user/index", fiber.Map{
		"Users": users,
	})
}
