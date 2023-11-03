package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/controller"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/controller/admin"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/controller/auth"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/controller/pemohon"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/middleware"
)

func SetupRoutes(app *fiber.App) {
	AdminAuth := middleware.AuthHandler{Role: "Admin"}
	PemohonAuth := middleware.AuthHandler{Role: "Pemohon"}

	app.Get("/", auth.Login)
	app.Post("/", auth.CheckLogin)
	app.Get("/logout", auth.Logout)
	app.Get("/register", middleware.GuestMiddleware, auth.Register)
	app.Post("/register", middleware.GuestMiddleware, auth.CheckRegister)

	adminGroup := app.Group("/admin", AdminAuth.AuthMiddleware)
	adminGroup.Get("/dashboard", admin.Dashboard)
	adminGroup.Get("/profile", admin.ProfileAdmin)
	adminGroup.Post("/profile", admin.UpdateAdmin)

	adminUser := adminGroup.Group("user")
	adminUser.Get("/list-admin", admin.ListAdmin)
	adminUser.Get("/add-admin", admin.AddAdmin)
	adminUser.Post("/add-admin", admin.StoreAdmin)
	adminUser.Get("/edit-admin/:userId", admin.EditAdmin)
	adminUser.Post("/edit-admin/:userId", admin.UpdateUserAdmin)
	adminUser.Get("/delete-admin/:userId", admin.DeleteAdmin)

	adminUser.Get("/list-pemohon", admin.ListPemohon)
	adminUser.Get("/delete-pemohon/:userId", admin.DeletePemohon)

	adminBrand := adminGroup.Group("brand")
	adminBrand.Get("/", admin.ListBrand)
	adminBrand.Get("/review/:brandId", admin.ReviewBrand)
	adminBrand.Post("/review/:brandId", admin.UpdateReviewBrand)

	pemohonGroup := app.Group("/pemohon", PemohonAuth.AuthMiddleware)
	pemohonGroup.Get("/dashboard", pemohon.Dashboard)

	pemohonProfile := pemohonGroup.Group("/profile")
	pemohonProfile.Get("user", pemohon.ProfilePemohon)
	pemohonProfile.Post("user", pemohon.UpdatePemohon)

	pemohonProfile.Get("business", pemohon.ProfileBusiness)
	pemohonProfile.Post("business", pemohon.UpdateBusiness)

	pemohonBrand := pemohonGroup.Group("brand") // tambahin middleware cek data bisnis sudah dilengkapi
	pemohonBrand.Get("/", pemohon.ListBrand)
	pemohonBrand.Get("/add", pemohon.AddBrand)
	pemohonBrand.Post("/add", pemohon.CreateBrand)
	pemohonBrand.Get("/detail/:brandId", pemohon.DetailBrand)
	pemohonBrand.Get("/edit/:brandId", pemohon.EditBrand) // tambahin middleware cek id edit = createdBy // sudah tpi msih di dalam function
	pemohonBrand.Post("/edit/:brandId", pemohon.UpdateBrand)
	pemohonBrand.Get("/delete/:brandId", pemohon.DeleteBrand)

	user := app.Group("/user")
	user.Get("/", controller.GetAllUsers)
	//user.Get("/create", controller.CreateUser)
	//user.Post("/create", controller.StoreUser)
}
