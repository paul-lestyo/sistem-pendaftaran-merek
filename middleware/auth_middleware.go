package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/paul-lestyo/sistem-pendaftaran-merek/helper"
	"strings"
)

type AuthHandler struct {
	Role string
}

func (ah *AuthHandler) AuthMiddleware(c *fiber.Ctx) error {
	if helper.GetSession(c, "RoleUser") == ah.Role {
		return c.Next()
	}
	return c.Redirect("/login")
}

func GuestMiddleware(c *fiber.Ctx) error {
	if helper.GetSession(c, "LoggedIn") != nil {
		role := helper.GetSession(c, "RoleUser").(string)
		return c.Redirect(strings.ToLower(role) + "/dashboard")
	}
	return c.Next()
}
