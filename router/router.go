package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/controller"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/controller/admin"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/controller/api"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/controller/auth"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/controller/pemohon"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/middleware"
)

func SetupRoutes(app *fiber.App) {
	AdminAuth := middleware.AuthHandler{Role: "Admin"}
	PemohonAuth := middleware.AuthHandler{Role: "Pemohon"}

	app.Get("/api/searchPDKI/:search", api.SearchDataPDKI)
	app.Get("/api/getChartBrand/:date", api.GetDataChartPermohonanMerek)
	app.Get("/api/getChartLogin/:date", api.GetDataChartLogin)

	app.Get("/", auth.Login)
	app.Post("/", auth.CheckLogin)
	app.Get("/logout", auth.Logout)
	app.Get("/register", middleware.GuestMiddleware, auth.Register)
	app.Post("/register", middleware.GuestMiddleware, auth.CheckRegister)

	adminGroup := app.Group("/admin", AdminAuth.AuthMiddleware)
	adminGroup.Get("/dashboard", admin.Dashboard)
	adminGroup.Get("/profile", admin.ProfileAdmin)
	adminGroup.Post("/profile", admin.UpdateProfileAdmin)

	adminUser := adminGroup.Group("user")
	adminUser.Get("/list-admin", admin.ListAdmin)
	adminUser.Get("/add-admin", admin.AddAdmin)
	adminUser.Post("/add-admin", admin.StoreAdmin)
	adminUser.Get("/edit-admin/:userId", admin.EditAdmin)
	adminUser.Post("/edit-admin/:userId", admin.UpdateAdmin)
	adminUser.Get("/delete-admin/:userId", admin.DeleteAdmin)

	adminUser.Get("/list-pemohon", admin.ListPemohon)
	adminUser.Get("/add-pemohon", admin.AddPemohon)
	adminUser.Post("/add-pemohon", admin.StorePemohon)
	adminUser.Get("/edit-pemohon/:userId", admin.EditPemohon)
	adminUser.Post("/edit-pemohon/:userId", admin.UpdatePemohon)
	adminUser.Get("/delete-pemohon/:userId", admin.DeletePemohon)
	adminUser.Get("deactivate-pemohon/:userId", admin.DeactivatePemohon)
	adminUser.Get("activate-pemohon/:userId", admin.ActivatePemohon)

	adminBrand := adminGroup.Group("brand")
	adminBrand.Get("/", admin.ListBrand)
	adminBrand.Get("/review/:brandId", admin.ReviewBrand)
	adminBrand.Post("/review/:brandId", admin.UpdateReviewBrand)

	adminAnnouncement := adminGroup.Group("announcement")
	adminAnnouncement.Get("/", admin.ListAnnouncement)
	adminAnnouncement.Get("/add", admin.AddAnnouncement)
	adminAnnouncement.Post("/add", admin.StoreAnnouncement)
	adminAnnouncement.Get("/edit/:announcementId", admin.EditAnnouncement)
	adminAnnouncement.Post("/edit/:announcementId", admin.UpdateAnnouncement)
	adminAnnouncement.Get("/delete/:announcementId", admin.DeleteAnnouncement)

	pemohonGroup := app.Group("/pemohon", PemohonAuth.AuthMiddleware)
	pemohonGroup.Get("/dashboard", pemohon.Dashboard)

	pemohonProfile := pemohonGroup.Group("/profile")
	pemohonProfile.Get("user", pemohon.ProfilePemohon)
	pemohonProfile.Post("user", pemohon.UpdatePemohon)

	pemohonProfile.Get("business", pemohon.ProfileBusiness)
	pemohonProfile.Post("business", pemohon.UpdateBusiness)

	pemohonBrand := pemohonGroup.Group("brand", middleware.BusinessFilledMiddleware) // tambahin middleware cek data bisnis sudah dilengkapi
	pemohonBrand.Get("/", pemohon.ListBrand)
	pemohonBrand.Get("/add", pemohon.AddBrand)
	pemohonBrand.Post("/add", pemohon.CreateBrand)
	pemohonBrand.Get("/detail/:brandId", pemohon.DetailBrand)
	pemohonBrand.Get("/edit/:brandId", pemohon.EditBrand) // tambahin middleware cek id edit = createdBy // sudah tpi msih di dalam function
	pemohonBrand.Post("/edit/:brandId", pemohon.UpdateBrand)
	pemohonBrand.Get("/delete/:brandId", pemohon.DeleteBrand)

	pemohonAnnouncement := pemohonGroup.Group("announcement")
	pemohonAnnouncement.Get("/", pemohon.ListAnnouncement)
	pemohonAnnouncement.Get("/:announcementId", pemohon.DetailAnnouncement)

	user := app.Group("/user")
	user.Get("/", controller.GetAllUsers)
	//user.Get("/create", controller.CreateUser)
	//user.Post("/create", controller.StoreUser)
}
