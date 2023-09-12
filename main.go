package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/database"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/router"
)

func main() {
	database.Connect()

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	router.SetupRoutes(app)

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	app.Listen(":8080")
}
