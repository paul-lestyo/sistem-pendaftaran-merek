package pemohon

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
	"os"
	"path"
)

func ListBrand(c *fiber.Ctx) error {
	message := helper.GetSession(c, "successMessage")
	helper.DeleteSession(c, "successMessage")
	var business model.Business
	var user model.User

	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	err = database.DB.Preload("Brands").First(&business, "user_id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	return c.Render("pemohon/brand/index", fiber.Map{
		"User":    user,
		"Brands":  business.Brands,
		"message": message,
	}, "layouts/pemohon")
}

func AddBrand(c *fiber.Ctx) error {
	message := helper.GetSession(c, "alertMessage")
	helper.DeleteSession(c, "alertMessage")
	var business model.Business
	var user model.User

	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	err = database.DB.Preload("Brands").First(&business, "user_id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	return c.Render("pemohon/brand/create", fiber.Map{
		"User":    user,
		"Brands":  business.Brands,
		"message": message,
	}, "layouts/pemohon")
}

type CreateUpdateBrandUser struct {
	BrandName string           `validate:"required,min=5,max=50"`
	DescBrand string           `validate:"required,min=5,max=50"`
	BrandLogo helper.FileInput `validate:"required,image_upload"`
}

func CreateBrand(c *fiber.Ctx) error {
	img, uploadImg := helper.CheckInputFile(c, "brand_logo")
	createBrandUser := CreateUpdateBrandUser{
		BrandName: c.FormValue("brand_name"),
		DescBrand: c.FormValue("desc_brand"),
		BrandLogo: img,
	}

	createBrandValidator := &helper.Validator{
		Validator: validate,
	}

	if errs := createBrandValidator.Validate(createBrandUser); len(errs) > 0 {
		return showCreateBrandErrors(c, createBrandUser, errs)
	}

	var business model.Business
	err := database.DB.First(&business, "user_id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	pathImg := ""
	if uploadImg {
		if path, ok := helper.UploadFile(c, "brand_logo", "brand"); ok {
			pathImg = path
		}
	}

	brand := model.Brand{
		BusinessID:  business.ID,
		BrandName:   c.FormValue("brand_name"),
		DescBrand:   c.FormValue("desc_brand"),
		BrandLogo:   pathImg,
		Status:      "Menunggu",
		Note:        "",
		CreatedByID: business.UserID,
		UpdatedByID: business.UserID,
	}

	err = database.DB.Create(&brand).Error
	helper.PanicIfError(err)

	helper.SetSession(c, "successMessage", "Permohonan Brand Berhasil Ditambahkan!")
	return c.Redirect("/pemohon/brand/")
}

func DetailBrand(c *fiber.Ctx) error {
	id := c.Params("brandId")
	var brand model.Brand
	var user model.User

	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	err = database.DB.First(&brand, "id = ?", id).Error
	helper.PanicIfError(err)

	if brand.CreatedByID.String() == helper.GetSession(c, "LoggedIn").(string) {
		return c.Render("pemohon/brand/detail", fiber.Map{
			"User":  user,
			"Brand": brand,
		}, "layouts/pemohon")
	} else {
		return c.SendStatus(404)
	}

}

func EditBrand(c *fiber.Ctx) error {
	message := helper.GetSession(c, "alertMessage")
	helper.DeleteSession(c, "alertMessage")

	id := c.Params("brandId")
	var brand model.Brand
	var user model.User

	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	err = database.DB.First(&brand, "id = ?", id).Error
	helper.PanicIfError(err)

	if brand.CreatedByID.String() == helper.GetSession(c, "LoggedIn").(string) {
		return c.Render("pemohon/brand/edit", fiber.Map{
			"User":    user,
			"Brand":   brand,
			"message": message,
		}, "layouts/pemohon")
	} else {
		return c.SendStatus(404)
	}

}

func UpdateBrand(c *fiber.Ctx) error {
	id := c.Params("brandId")
	var brand model.Brand

	err := database.DB.First(&brand, "id = ?", id).Error
	helper.PanicIfError(err)

	if brand.CreatedByID.String() != helper.GetSession(c, "LoggedIn").(string) {
		return c.SendStatus(404)
	}

	img, uploadImg := helper.CheckInputFile(c, "brand_logo")
	updateBrandUser := CreateUpdateBrandUser{
		BrandName: c.FormValue("brand_name"),
		DescBrand: c.FormValue("desc_brand"),
		BrandLogo: img,
	}

	updateBrandValidator := &helper.Validator{
		Validator: validate,
	}

	if errs := updateBrandValidator.Validate(updateBrandUser); len(errs) > 0 {
		return showUpdateBrandErrors(c, brand, updateBrandUser, errs)
	}

	var business model.Business
	err = database.DB.First(&business, "user_id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	brand.BrandName = c.FormValue("brand_name")
	brand.DescBrand = c.FormValue("desc_brand")

	if uploadImg {
		if pathFile, ok := helper.UploadFile(c, "brand_logo", "brand"); ok {
			_ = os.Remove(path.Join("assets", brand.BrandLogo))
			brand.BrandLogo = pathFile
		}
	}

	err = database.DB.Save(&brand).Error
	helper.PanicIfError(err)

	helper.SetSession(c, "successMessage", "Permohonan Brand Berhasil Diedit!")
	return c.Redirect("/pemohon/brand/")
}

func DeleteBrand(c *fiber.Ctx) error {
	id := c.Params("brandId")
	var brand model.Brand

	err := database.DB.First(&brand, "id = ?", id).Error
	helper.PanicIfError(err)

	if brand.CreatedByID.String() != helper.GetSession(c, "LoggedIn").(string) {
		return c.SendStatus(404)
	}

	_ = os.Remove(path.Join("assets", brand.BrandLogo))
	database.DB.Delete(&brand)
	return c.Redirect("/pemohon/brand/")
}

type MessageCreateUpdateBrand struct {
	BrandName string
	DescBrand string
	BrandLogo string
}

func showCreateBrandErrors(c *fiber.Ctx, oldInput CreateUpdateBrandUser, errs map[string]string) error {
	var user model.User

	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)
	var errsStruct = MessageCreateUpdateBrand{}
	if err := mapstructure.Decode(errs, &errsStruct); err != nil {
		panic(err)
	}
	return c.Render("pemohon/brand/create", fiber.Map{
		"User":     user,
		"oldInput": oldInput,
		"Errors":   errsStruct,
	}, "layouts/pemohon")
}

func showUpdateBrandErrors(c *fiber.Ctx, brand model.Brand, oldInput CreateUpdateBrandUser, errs map[string]string) error {
	var user model.User

	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)
	var errsStruct = MessageCreateUpdateBrand{}
	if err := mapstructure.Decode(errs, &errsStruct); err != nil {
		panic(err)
	}
	return c.Render("pemohon/brand/edit", fiber.Map{
		"User":     user,
		"Brand":    brand,
		"oldInput": oldInput,
		"Errors":   errsStruct,
	}, "layouts/pemohon")
}
