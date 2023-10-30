package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/controller"
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

	admin := app.Group("/admin", AdminAuth.AuthMiddleware)
	admin.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.SendString("I'm a GET request dashboard admin!")
	})

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
