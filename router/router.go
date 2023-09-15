package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/controller"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/controller/auth"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/login", auth.Login)
	app.Post("login", auth.CheckLogin)
	app.Get("/register", auth.Register)
	app.Post("/register", auth.CheckRegister)

	user := app.Group("/user")

	user.Get("/", controller.GetAllUsers)
	//user.Get("/create", controller.CreateUser)
	//user.Post("/create", controller.StoreUser)
}
