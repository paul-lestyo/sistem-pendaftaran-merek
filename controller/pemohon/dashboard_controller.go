package pemohon

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
)

type Infografis struct {
	TotalPermohonanMerek    int64
	TotalPermohonanPerbaiki int64
	TotalPermohonanOK       int64
	TotalLogin              int64
}

func Dashboard(c *fiber.Ctx) error {
	var user model.User
	var infografis Infografis

	err := database.DB.Preload("Business").First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	database.DB.Model(&model.Brand{}).Where("business_id = ?", user.Business.ID).Count(&infografis.TotalPermohonanMerek)
	database.DB.Model(&model.Brand{}).Where("business_id = ?", user.Business.ID).Where("status = ?", "OK").Count(&infografis.TotalPermohonanOK)
	database.DB.Model(&model.Brand{}).Where("business_id = ?", user.Business.ID).Where("status = ?", "Perbaiki").Count(&infografis.TotalPermohonanPerbaiki)
	database.DB.Model(&model.Log{}).Where("user_id", user.ID).Count(&infografis.TotalLogin)
	fmt.Println(infografis.TotalPermohonanMerek)

	return c.Render("pemohon/dashboard", fiber.Map{
		"User":       user,
		"Infografis": infografis,
	}, "layouts/pemohon")
}
