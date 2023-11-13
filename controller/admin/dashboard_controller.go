package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
)

type Infografis struct {
	TotalAkun            int64
	TotalPermohonanMerek int64
	TotalPengumuman      int64
	TotalLogin           int64
}

func Dashboard(c *fiber.Ctx) error {
	var user model.User
	var infografis Infografis

	database.DB.Model(&model.User{}).Count(&infografis.TotalAkun)
	database.DB.Model(&model.Brand{}).Count(&infografis.TotalPermohonanMerek)
	database.DB.Model(&model.Announcement{}).Count(&infografis.TotalPengumuman)
	database.DB.Model(&model.Log{}).Count(&infografis.TotalLogin)

	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)
	return c.Render("admin/dashboard", fiber.Map{
		"User":       user,
		"Infografis": infografis,
	}, "layouts/admin")
}
