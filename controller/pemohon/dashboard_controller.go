package pemohon

import (
	"github.com/gofiber/fiber/v2"
)

func Dashboard(c *fiber.Ctx) error {
	return c.Render("pemohon/dashboard", fiber.Map{}, "layouts/pemohon")
}
