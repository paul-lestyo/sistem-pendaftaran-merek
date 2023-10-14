package pemohon

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
)

func ListBrand(c *fiber.Ctx) error {
	message := helper.GetSession(c, "alertMessage")
	helper.DeleteSession(c, "alertMessage")
	var business model.Business

	err := database.DB.Preload("Brands").First(&business, "user_id = ?", helper.GetSession(c, "LoggedIn")).Error
	fmt.Println(business.Brands)
	helper.PanicIfError(err)

	return c.Render("pemohon/brand/index", fiber.Map{
		"Brands":  business.Brands,
		"message": message,
	}, "layouts/pemohon")
}

func AddBrand(c *fiber.Ctx) error {
	message := helper.GetSession(c, "alertMessage")
	helper.DeleteSession(c, "alertMessage")
	var business model.Business

	err := database.DB.Preload("Brands").First(&business, "user_id = ?", helper.GetSession(c, "LoggedIn")).Error
	fmt.Println(business.Brands)
	helper.PanicIfError(err)

	return c.Render("pemohon/brand/create", fiber.Map{
		"Brands":  business.Brands,
		"message": message,
	}, "layouts/pemohon")
}

func EditBrand(c *fiber.Ctx) error {
	message := helper.GetSession(c, "alertMessage")
	helper.DeleteSession(c, "alertMessage")

	id := c.Params("brandId")
	var brand model.Brand

	err := database.DB.First(&brand, "id = ?", id).Error
	fmt.Println(brand)
	helper.PanicIfError(err)

	return c.Render("pemohon/brand/edit", fiber.Map{
		"Brand":   brand,
		"message": message,
	}, "layouts/pemohon")
}
