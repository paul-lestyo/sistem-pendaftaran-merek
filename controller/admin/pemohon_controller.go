package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
)

func ListPemohon(c *fiber.Ctx) error {
	message := helper.GetSession(c, "successMessage")
	helper.DeleteSession(c, "successMessage")
	var user model.User

	var role model.Role
	err := database.DB.First(&role, "name = ?", "Pemohon").Error
	helper.PanicIfError(err)

	err = database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	rows, err := database.DB.Raw("SELECT users.id, users.name, users.email, users.image_url, businesses.business_name, businesses.owner_name, COUNT(brands.id) as count_brands FROM users "+
		"left join businesses on businesses.user_id = users.id "+
		"left join brands on businesses.id = brands.business_id "+
		"WHERE users.role_id=? GROUP BY users.id;", role.ID).Rows()
	defer rows.Close()
	helper.PanicIfError(err)

	var data []model.UsersBusinessCountBrands
	for rows.Next() {
		database.DB.ScanRows(rows, &data)
	}

	return c.Render("admin/user/list-pemohon", fiber.Map{
		"User":    user,
		"Users":   data,
		"message": message,
	}, "layouts/admin")
}

func DeletePemohon(c *fiber.Ctx) error {
	id := c.Params("userId")
	var user model.User

	err := database.DB.First(&user, "id = ?", id).Error
	helper.PanicIfError(err)

	//_ = os.Remove(path.Join("assets", user.BrandLogo))
	err = database.DB.Delete(&user).Error
	helper.PanicIfError(err)

	helper.SetSession(c, "successMessage", "Berhasil menghapus akun Pemohon!")
	return c.Redirect("/admin/user/list-pemohon")
}
