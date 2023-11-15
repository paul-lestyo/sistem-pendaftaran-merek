package guest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
)

func Dashboard(c *fiber.Ctx) error {
	var brands []model.Brand

	err := database.DB.Find(&brands).Error
	helper.PanicIfError(err)

	return c.Render("guest/dashboard", fiber.Map{
		"Brands": brands,
	}, "layouts/guest")
}
