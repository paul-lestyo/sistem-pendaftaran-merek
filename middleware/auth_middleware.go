package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
)

func GuestMiddleware(c *fiber.Ctx) error {
	if helper.GetSession(c, "LoggedIn") != nil {
		return c.Redirect("/dashboard")
	}
	return c.Next()
}
