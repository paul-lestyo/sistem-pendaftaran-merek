package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
)

func BusinessFilledMiddleware(c *fiber.Ctx) error {
	var user model.User

	err := database.DB.Preload("Business").First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	if user.Business.BusinessName != "" && user.Business.BusinessAddress != "" {
		return c.Next()
	}

	helper.SetSession(c, "messageAlert", "Silahkan Lengkapi Data Bisnis Anda untuk dapat Mengajukan Pendaftaran Merek!")
	return c.Redirect("/pemohon/profile/business")

}
