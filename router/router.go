package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/handler"
)

func SetupRoutes(app *fiber.App) {
	user := app.Group("/user")

	user.Get("/", handler.GetAllUsers)
	//user.Get("/create", handler.CreateUser)
	//user.Post("/create", handler.StoreUser)
}
