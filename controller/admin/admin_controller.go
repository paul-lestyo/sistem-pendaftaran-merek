package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
)

func ListAdmin(c *fiber.Ctx) error {
	message := helper.GetSession(c, "successMessage")
	helper.DeleteSession(c, "successMessage")
	var users []model.User
	var user model.User

	var role model.Role
	err := database.DB.First(&role, "name = ?", "Admin").Error
	helper.PanicIfError(err)

	err = database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	err = database.DB.Where("role_id = ?", role.ID).Find(&users).Error
	helper.PanicIfError(err)

	return c.Render("admin/user/list-admin", fiber.Map{
		"User":    user,
		"Users":   users,
		"message": message,
	}, "layouts/admin")
}
