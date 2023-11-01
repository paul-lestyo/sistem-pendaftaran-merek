package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
)

func ListBrand(c *fiber.Ctx) error {
	message := helper.GetSession(c, "successMessage")
	helper.DeleteSession(c, "successMessage")
	var brands []model.Brand
	var user model.User

	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	err = database.DB.Find(&brands).Error
	helper.PanicIfError(err)

	return c.Render("admin/brand/index", fiber.Map{
		"User":    user,
		"Brands":  brands,
		"message": message,
	}, "layouts/admin")
}

func ReviewBrand(c *fiber.Ctx) error {
	id := c.Params("brandId")
	var brand model.Brand
	var user model.User

	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	err = database.DB.Preload("Business").First(&brand, "id = ?", id).Error
	helper.PanicIfError(err)

	return c.Render("admin/brand/review", fiber.Map{
		"User":     user,
		"Business": brand.Business,
		"Brand":    brand,
	}, "layouts/admin")
}

type UpdateReviewBrandPemohon struct {
	Status string `validate:"required"`
	Note   string `validate:"omitempty"`
}

func UpdateReviewBrand(c *fiber.Ctx) error {
	id := c.Params("brandId")
	var brand model.Brand

	err := database.DB.Preload("Business").First(&brand, "id = ?", id).Error
	helper.PanicIfError(err)

	updateReviewBrandPemohon := UpdateReviewBrandPemohon{
		Status: c.FormValue("status"),
		Note:   c.FormValue("note"),
	}

	updateBrandValidator := &helper.Validator{
		Validator: validate,
	}

	if errs := updateBrandValidator.Validate(updateReviewBrandPemohon); len(errs) > 0 {
		return showUpdateReviewBrandErrors(c, brand, updateReviewBrandPemohon, errs)
	}

	brand.Status = c.FormValue("status")
	brand.Note = c.FormValue("note")

	err = database.DB.Save(&brand).Error
	helper.PanicIfError(err)

	helper.SetSession(c, "successMessage", "Review Permohonan Brand telah Berhasil!")
	return c.Redirect("/admin/brand/")
}

type MessageUpdateReviewBrand struct {
	Status string
	Note   string
}

func showUpdateReviewBrandErrors(c *fiber.Ctx, brand model.Brand, oldInput UpdateReviewBrandPemohon, errs map[string]string) error {
	var user model.User

	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)
	var errsStruct = MessageUpdateReviewBrand{}
	if err := mapstructure.Decode(errs, &errsStruct); err != nil {
		panic(err)
	}
	return c.Render("admin/brand/review", fiber.Map{
		"User":     user,
		"Brand":    brand,
		"Business": brand.Business,
		"oldInput": oldInput,
		"Errors":   errsStruct,
	}, "layouts/admin")
}
