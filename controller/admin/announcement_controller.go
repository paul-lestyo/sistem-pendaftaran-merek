package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
)

func ListAnnouncement(c *fiber.Ctx) error {
	message := helper.GetSession(c, "successMessage")
	helper.DeleteSession(c, "successMessage")
	var announcements []model.Announcement
	var user model.User

	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	err = database.DB.Find(&announcements).Error
	helper.PanicIfError(err)

	return c.Render("admin/announcement/index", fiber.Map{
		"User":          user,
		"Announcements": announcements,
		"message":       message,
	}, "layouts/admin")
}

func AddAnnouncement(c *fiber.Ctx) error {
	var user model.User

	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)
	return c.Render("admin/announcement/add", fiber.Map{
		"User": user,
	}, "layouts/admin")
}

type AddEditAnnouncementData struct {
	Title    string           `validate:"required,min=5,max=255"`
	Desc     string           `validate:"required,min=5"`
	Tag      string           `validate:"required,min=3"`
	ImageUrl helper.FileInput `validate:"required,image_upload"`
}

func StoreAnnouncement(c *fiber.Ctx) error {
	imgAnnouncement, updateImgAnnouncement := helper.CheckInputFile(c, "image_url")

	var user model.User
	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	registerUser := AddEditAnnouncementData{
		Title:    c.FormValue("title"),
		Desc:     c.FormValue("desc"),
		Tag:      c.FormValue("tag"),
		ImageUrl: imgAnnouncement,
	}

	registerValidator := &helper.Validator{
		Validator: validate,
	}

	if errs := registerValidator.Validate(registerUser); len(errs) > 0 {
		return showAddAnnouncementErrors(c, registerUser, errs)
	}

	path := ""
	if updateImgAnnouncement {
		if pathLogo, ok := helper.UploadFile(c, "image_url", "announcement/"); ok {
			path = pathLogo
		}
	}

	announcement := model.Announcement{
		Title:     c.FormValue("title"),
		Desc:      c.FormValue("desc"),
		Tag:       c.FormValue("tag"),
		ImageUrl:  path,
		CreatedBy: user.ID,
	}

	result := database.DB.Create(&announcement)

	if result != nil {
		helper.SetSession(c, "successMessage", "Berhasil menambahkan Pengumuman!")
		return c.Redirect("/admin/announcement/")
	}
	return nil
}

func EditAnnouncement(c *fiber.Ctx) error {
	id := c.Params("announcementId")
	var user model.User
	var announcement model.Announcement

	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	err = database.DB.First(&announcement, "id = ?", id).Error
	helper.PanicIfError(err)

	return c.Render("admin/announcement/edit", fiber.Map{
		"User":         user,
		"Announcement": announcement,
	}, "layouts/admin")
}

func UpdateAnnouncement(c *fiber.Ctx) error {
	id := c.Params("announcementId")
	var announcement model.Announcement

	err := database.DB.First(&announcement, "id = ?", id).Error
	helper.PanicIfError(err)

	imgAnnouncement, updateImgAnnouncement := helper.CheckInputFile(c, "image_url")

	annountcementEdit := AddEditAnnouncementData{
		Title:    c.FormValue("title"),
		Desc:     c.FormValue("desc"),
		Tag:      c.FormValue("tag"),
		ImageUrl: imgAnnouncement,
	}

	registerValidator := &helper.Validator{
		Validator: validate,
	}

	if errs := registerValidator.Validate(annountcementEdit); len(errs) > 0 {
		return showEditAnnouncementErrors(c, annountcementEdit, errs)
	}

	announcement.Title = c.FormValue("title")
	announcement.Desc = c.FormValue("desc")
	announcement.Tag = c.FormValue("tag")

	if updateImgAnnouncement {
		if path, ok := helper.UploadFile(c, "image_url", "announcement/"); ok {
			announcement.ImageUrl = path
		}
	}

	err = database.DB.Save(&announcement).Error
	helper.PanicIfError(err)

	helper.SetSession(c, "successMessage", "Berhasil mengedit Pengumuman!")
	return c.Redirect("/admin/announcement/")
}

func DeleteAnnouncement(c *fiber.Ctx) error {
	id := c.Params("announcementId")
	var announcement model.Announcement

	err := database.DB.First(&announcement, "id = ?", id).Error
	helper.PanicIfError(err)

	//_ = os.Remove(path.Join("assets", user.BrandLogo))
	err = database.DB.Delete(&announcement).Error
	helper.PanicIfError(err)

	helper.SetSession(c, "successMessage", "Berhasil menghapus Pengumuman!")
	return c.Redirect("/admin/announcement/")
}

type MessageAddEditAnnouncement struct {
	Title    string
	Desc     string
	Tag      string
	ImageUrl string
}

func showAddAnnouncementErrors(c *fiber.Ctx, oldInput AddEditAnnouncementData, errs map[string]string) error {
	var errsStruct = MessageAddEditAnnouncement{}
	if err := mapstructure.Decode(errs, &errsStruct); err != nil {
		panic(err)
	}
	return c.Render("admin/announcement/add", fiber.Map{
		"oldInput": oldInput,
		"Errors":   errsStruct,
	}, "layouts/admin")
}

func showEditAnnouncementErrors(c *fiber.Ctx, oldInput AddEditAnnouncementData, errs map[string]string) error {
	var errsStruct = MessageAddEditAnnouncement{}
	if err := mapstructure.Decode(errs, &errsStruct); err != nil {
		panic(err)
	}
	return c.Render("admin/announcement/edit", fiber.Map{
		"oldInput": oldInput,
		"Errors":   errsStruct,
	}, "layouts/admin")
}
