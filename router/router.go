package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/controller"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/controller/auth"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/middleware"
)

func SetupRoutes(app *fiber.App) {
	AdminAuth := middleware.AuthHandler{Role: "Admin"}
	PemohonAuth := middleware.AuthHandler{Role: "Pemohon"}

	app.Get("/login", auth.Login)
	app.Post("login", auth.CheckLogin)
	app.Get("/register", middleware.GuestMiddleware, auth.Register)
	app.Post("/register", middleware.GuestMiddleware, auth.CheckRegister)

	admin := app.Group("/admin")
	admin.Get("/dashboard", AdminAuth.AuthMiddleware, func(c *fiber.Ctx) error {
		return c.SendString("I'm a GET request dashboard admin!")
	})

	pemohon := app.Group("/pemohon")
	pemohon.Get("/dashboard", PemohonAuth.AuthMiddleware, func(c *fiber.Ctx) error {
		return c.SendString("I'm a GET request dashboard pemohon!")
	})

	user := app.Group("/user")

	user.Get("/", controller.GetAllUsers)
	//user.Get("/create", controller.CreateUser)
	//user.Post("/create", controller.StoreUser)
}
