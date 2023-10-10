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

	app.Get("/login", auth.Login)
	app.Post("login", auth.CheckLogin)
	app.Get("/register", middleware.GuestMiddleware, auth.Register)
	app.Post("/register", middleware.GuestMiddleware, auth.CheckRegister)

	admin := app.Group("/admin", AdminAuth.AuthMiddleware)
	admin.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.SendString("I'm a GET request dashboard admin!")
	})

	pemohonGroup := app.Group("/pemohon", PemohonAuth.AuthMiddleware)
	pemohonGroup.Get("/dashboard", pemohon.Dashboard)

	pemohonGroup.Get("profile", pemohon.ProfilePemohon)
	pemohonGroup.Post("profile", pemohon.UpdatePemohon)

	user := app.Group("/user")

	user.Get("/", controller.GetAllUsers)
	//user.Get("/create", controller.CreateUser)
	//user.Post("/create", controller.StoreUser)
}
