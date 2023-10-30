package pemohon

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/model"
)

type UpdateBusinessVal struct {
	BusinessName      string           `validate:"required,min=5,max=50" name:"Nama Bisnis"`
	BusinessAddress   string           `validate:"required,min=5,max=50" name:"Alamat Bisnis"`
	OwnerName         string           `validate:"required,min=5,max=50" name:"Nama Owner"`
	BusinessLogo      helper.FileInput `validate:"required,image_upload" name:"Logo Bisnis"`
	UMKCertificateUrl helper.FileInput `validate:"omitempty,image_upload" name:"Surat Keterangan UMK"`
	SignatureUrl      helper.FileInput `validate:"required,image_upload" name:"Tanda Tangan"`
}

func ProfileBusiness(c *fiber.Ctx) error {
	message := helper.GetSession(c, "successMessage")
	helper.DeleteSession(c, "successMessage")
	var user model.User

	err := database.DB.Preload("Business").First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	return c.Render("pemohon/profile/business", fiber.Map{
		"User":     user,
		"Business": user.Business,
		"message":  message,
	}, "layouts/pemohon")
}

func UpdateBusiness(c *fiber.Ctx) error {
	imgCertificate, updateImgCertificate := helper.CheckInputFile(c, "umk_certificate_url")
	imgSignature, updateImgSignature := helper.CheckInputFile(c, "signature_url")
	imgLogo, updateImgLogo := helper.CheckInputFile(c, "business_logo")

	updateBusinessVal := UpdateBusinessVal{
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

	if errs := registerValidator.Validate(updateBusinessVal); len(errs) > 0 {
		return showProfileBusinessErrors(c, updateBusinessVal, errs)
	}

	var business model.Business
	err := database.DB.First(&business, "user_id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)

	business.BusinessName = c.FormValue("business_name")
	business.BusinessAddress = c.FormValue("business_address")
	business.OwnerName = c.FormValue("owner_name")
	if updateImgLogo {
		if pathLogo, ok := helper.UploadFile(c, "business_logo", "profile/business"); ok {
			business.BusinessLogo = pathLogo
		}
	}

	if updateImgCertificate {
		if pathCertificate, ok := helper.UploadFile(c, "umk_certificate_url", "profile/business"); ok {
			business.UMKCertificateUrl = pathCertificate
		}
	}

	if updateImgSignature {
		if pathSignature, ok := helper.UploadFile(c, "signature_url", "profile/business"); ok {
			business.SignatureUrl = pathSignature
		}
	}

	err = database.DB.Save(&business).Error
	helper.PanicIfError(err)

	helper.SetSession(c, "successMessage", "Profile Bisnis Berhasil Tersimpan!")
	return c.Redirect("/pemohon/profile/business")
}

type MessageBusiness struct {
	BusinessName      string
	BusinessAddress   string
	OwnerName         string
	BusinessLogo      string
	UMKCertificateUrl string
	SignatureUrl      string
}

func showProfileBusinessErrors(c *fiber.Ctx, oldInput UpdateBusinessVal, errs map[string]string) error {
	var business model.Business
	var user model.User

	err := database.DB.First(&user, "id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)
	err = database.DB.First(&business, "user_id = ?", helper.GetSession(c, "LoggedIn")).Error
	helper.PanicIfError(err)
	var errsStruct = MessageBusiness{}
	if err := mapstructure.Decode(errs, &errsStruct); err != nil {
		panic(err)
	}
	return c.Render("pemohon/profile/business", fiber.Map{
		"User":     user,
		"Business": business,
		"oldInput": oldInput,
		"Errors":   errsStruct,
	}, "layouts/pemohon")
}
