package admin

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
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

func AddPemohon(c *fiber.Ctx) error {
	var user model.User

	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)
	return c.Render("admin/user/add-pemohon", fiber.Map{
		"User": user,
	}, "layouts/admin")
}

type AddPemohonUser struct {
	Name              string           `validate:"required,min=5,max=50"`
	Email             string           `validate:"required,min=5,email,unique email"`
	Password          string           `validate:"required,min=3"`
	ImageUrl          helper.FileInput `validate:"omitempty,image_upload"`
	BusinessName      string           `validate:"required,min=5,max=50" name:"Nama Bisnis"`
	BusinessAddress   string           `validate:"required,min=5,max=50" name:"Alamat Bisnis"`
	OwnerName         string           `validate:"required,min=5,max=50" name:"Nama Owner"`
	BusinessLogo      helper.FileInput `validate:"omitempty,image_upload" name:"Logo Bisnis"`
	UMKCertificateUrl helper.FileInput `validate:"omitempty,image_upload" name:"Surat Keterangan UMK"`
	SignatureUrl      helper.FileInput `validate:"omitempty,image_upload" name:"Tanda Tangan"`
}

func StorePemohon(c *fiber.Ctx) error {
	imgProfile, updateImgProfile := helper.CheckInputFile(c, "image_url")
	imgCertificate, updateImgCertificate := helper.CheckInputFile(c, "umk_certificate_url")
	imgSignature, updateImgSignature := helper.CheckInputFile(c, "signature_url")
	imgLogo, updateImgLogo := helper.CheckInputFile(c, "business_logo")

	var role model.Role
	err := database.DB.First(&role, "name = ?", "Pemohon").Error
	helper.PanicIfError(err)

	registerUser := AddPemohonUser{
		Name:              c.FormValue("name"),
		Email:             c.FormValue("email"),
		Password:          c.FormValue("password"),
		ImageUrl:          imgProfile,
		BusinessName:      c.FormValue("business_name"),
		BusinessAddress:   c.FormValue("business_address"),
		OwnerName:         c.FormValue("owner_name"),
		BusinessLogo:      imgLogo,
		UMKCertificateUrl: imgCertificate,
		SignatureUrl:      imgSignature,
	}

	registerValidator := &helper.Validator{
		Validator: validate,
	}

	if errs := registerValidator.Validate(registerUser); len(errs) > 0 {
		return showAddPemohonErrors(c, registerUser, errs)
	}

	pathProfile := ""
	pathLogo := ""
	pathCertificate := ""
	pathSignature := ""

	if updateImgProfile {
		if path, ok := helper.UploadFile(c, "image_url", "profile/user"); ok {
			pathProfile = path
		}
	}

	if updateImgLogo {
		if path, ok := helper.UploadFile(c, "business_logo", "profile/business"); ok {
			pathLogo = path
		}
	}

	if updateImgCertificate {
		if path, ok := helper.UploadFile(c, "umk_certificate_url", "profile/business"); ok {
			pathCertificate = path
		}
	}

	if updateImgSignature {
		if path, ok := helper.UploadFile(c, "signature_url", "profile/business"); ok {
			pathSignature = path
		}
	}

	hashedPassword, _ := helper.HashPassword(c.FormValue("password"))
	user := model.User{
		Name:     c.FormValue("name"),
		Email:    c.FormValue("email"),
		Password: hashedPassword,
		ImageUrl: pathProfile,
		RoleID:   role.ID,
		Business: &model.Business{
			BusinessName:      c.FormValue("business_name"),
			BusinessAddress:   c.FormValue("business_address"),
			OwnerName:         c.FormValue("owner_name"),
			BusinessLogo:      pathLogo,
			UMKCertificateUrl: pathCertificate,
			SignatureUrl:      pathSignature,
		},
	}

	result := database.DB.Create(&user)

	if result != nil {
		helper.SetSession(c, "successMessage", "Berhasil menambahkan akun Pemohon!")
		return c.Redirect("/admin/user/list-pemohon")
	}
	return nil
}

func EditPemohon(c *fiber.Ctx) error {
	id := c.Params("userId")
	var user model.User
	var userEdit model.User

	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	err = database.DB.Preload("Business").First(&userEdit, "id = ?", id).Error
	helper.PanicIfError(err)

	fmt.Println("hihi:", userEdit)

	return c.Render("admin/user/edit-pemohon", fiber.Map{
		"User":     user,
		"UserEdit": userEdit,
		"Business": userEdit.Business,
	}, "layouts/admin")
}

type EditPemohonUser struct {
	Name              string           `validate:"required,min=5,max=50"`
	Email             string           `validate:"required,min=5,email"`
	Password          string           `validate:"omitempty,min=3"`
	ImageUrl          helper.FileInput `validate:"omitempty,image_upload"`
	BusinessName      string           `validate:"required,min=5,max=50" name:"Nama Bisnis"`
	BusinessAddress   string           `validate:"required,min=5,max=50" name:"Alamat Bisnis"`
	OwnerName         string           `validate:"required,min=5,max=50" name:"Nama Owner"`
	BusinessLogo      helper.FileInput `validate:"omitempty,image_upload" name:"Logo Bisnis"`
	UMKCertificateUrl helper.FileInput `validate:"omitempty,image_upload" name:"Surat Keterangan UMK"`
	SignatureUrl      helper.FileInput `validate:"omitempty,image_upload" name:"Tanda Tangan"`
}

func UpdatePemohon(c *fiber.Ctx) error {
	id := c.Params("userId")
	var user model.User
	var business model.Business

	err := database.DB.First(&user, "id = ?", id).Error
	helper.PanicIfError(err)

	err = database.DB.First(&business, "user_id = ?", id).Error
	helper.PanicIfError(err)

	imgProfile, updateImgProfile := helper.CheckInputFile(c, "image_url")
	imgCertificate, updateImgCertificate := helper.CheckInputFile(c, "umk_certificate_url")
	imgSignature, updateImgSignature := helper.CheckInputFile(c, "signature_url")
	imgLogo, updateImgLogo := helper.CheckInputFile(c, "business_logo")

	editUser := EditPemohonUser{
		Name:              c.FormValue("name"),
		Email:             c.FormValue("email"),
		Password:          c.FormValue("password"),
		ImageUrl:          imgProfile,
		BusinessName:      c.FormValue("business_name"),
		BusinessAddress:   c.FormValue("business_address"),
		OwnerName:         c.FormValue("owner_name"),
		BusinessLogo:      imgLogo,
		UMKCertificateUrl: imgCertificate,
		SignatureUrl:      imgSignature,
	}

	registerValidator := &helper.Validator{
		Validator: validate,
	}

	if errs := registerValidator.Validate(editUser); len(errs) > 0 {
		return showEditPemohonErrors(c, editUser, errs)
	}

	user.Name = c.FormValue("name")
	user.Email = c.FormValue("email")
	if pass := c.FormValue("password"); pass != "" {
		hashedPassword, _ := helper.HashPassword(pass)
		user.Password = hashedPassword
	}

	if updateImgProfile {
		if path, ok := helper.UploadFile(c, "image_url", "profile/user"); ok {
			user.ImageUrl = path
		}
	}

	business.BusinessName = c.FormValue("business_name")
	business.BusinessAddress = c.FormValue("business_address")
	business.OwnerName = c.FormValue("owner_name")

	if updateImgLogo {
		if path, ok := helper.UploadFile(c, "business_logo", "profile/business"); ok {
			business.BusinessLogo = path
		}
	}

	if updateImgCertificate {
		if path, ok := helper.UploadFile(c, "umk_certificate_url", "profile/business"); ok {
			business.UMKCertificateUrl = path
		}
	}

	if updateImgSignature {
		if path, ok := helper.UploadFile(c, "signature_url", "profile/business"); ok {
			business.SignatureUrl = path
		}
	}
	err = database.DB.Save(&user).Error
	err = database.DB.Save(&business).Error
	helper.PanicIfError(err)

	helper.SetSession(c, "successMessage", "Berhasil mengedit akun Pemohon!")
	return c.Redirect("/admin/user/list-pemohon")
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

type MessageAddEditPemohon struct {
	Name              string
	Email             string
	Password          string
	ImageUrl          string
	BusinessName      string
	BusinessAddress   string
	OwnerName         string
	BusinessLogo      string
	UMKCertificateUrl string
	SignatureUrl      string
}

func showAddPemohonErrors(c *fiber.Ctx, oldInput AddPemohonUser, errs map[string]string) error {
	var errsStruct = MessageAddEditPemohon{}
	if err := mapstructure.Decode(errs, &errsStruct); err != nil {
		panic(err)
	}
	return c.Render("admin/user/add-pemohon", fiber.Map{
		"oldInput": oldInput,
		"Errors":   errsStruct,
	}, "layouts/admin")
}

func showEditPemohonErrors(c *fiber.Ctx, oldInput EditPemohonUser, errs map[string]string) error {
	var errsStruct = MessageAddEditPemohon{}
	if err := mapstructure.Decode(errs, &errsStruct); err != nil {
		panic(err)
	}
	return c.Render("admin/user/edit-pemohon", fiber.Map{
		"oldInput": oldInput,
		"Errors":   errsStruct,
	}, "layouts/admin")
}
