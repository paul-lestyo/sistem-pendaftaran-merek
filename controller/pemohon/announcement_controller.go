package pemohon

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
)

type Announcement struct {
	ID        uuid.UUID
	Title     string
	Desc      string
	Tag       string
	ImageUrl  string
	CreatedBy uuid.UUID
	CreatedAt string
}

func ListAnnouncement(c *fiber.Ctx) error {
	message := helper.GetSession(c, "successMessage")
	helper.DeleteSession(c, "successMessage")
	var announcements []Announcement
	var user model.User

	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	err = database.DB.Select("`id`, `title`, `desc`, `tag`, `image_url`, DATE_FORMAT(created_at, '%Y-%m-%d') as created_at").
		Find(&announcements).Error

	helper.PanicIfError(err)

	return c.Render("pemohon/announcement/index", fiber.Map{
		"User":          user,
		"Announcements": announcements,
		"message":       message,
	}, "layouts/pemohon")
}

func DetailAnnouncement(c *fiber.Ctx) error {
	id := c.Params("announcementId")
	var announcement Announcement
	var user model.User

	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	err = database.DB.Select("`id`, `title`, `desc`, `tag`, `image_url`, DATE_FORMAT(created_at, '%Y-%m-%d') as created_at").
		First(&announcement, "id = ?", id).Error
	helper.PanicIfError(err)

	return c.Render("pemohon/announcement/detail", fiber.Map{
		"User":         user,
		"Announcement": announcement,
	}, "layouts/pemohon")

}
